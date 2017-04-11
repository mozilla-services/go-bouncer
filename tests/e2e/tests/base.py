# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.

from urllib import urlencode
from urlparse import urlparse

import requests


class Base:

    cdn_netloc_locations = [
        'download.cdn.mozilla.net',
        'download-installer.cdn.mozilla.net'
    ]
    _user_agent_firefox = ('Mozilla/5.0 (Macintosh; Intel Mac OS X 10.7; '
                           'rv:10.0.1) Gecko/20100101 Firefox/10.0.1')

    def get_expected_fx_pkg_str(self, os, alias, product_version):
        """Output a string representation of the expected Firefox filename.

        :param os: string, 'win', 'osx', and 'linux' are currently supported.
        :param alias: string, the go-bouncer alias. Examples include firefox-latest,
        firefox-beta-latest, etc.
        :param product_version: string, the Firefox version number.

        :return: string
        """
        if os in ('win', 'win64'):
            if 'aurora' in alias or 'nightly' in alias:
                return 'firefox-{0}.en-US.win32.installer.exe'.format(product_version)
            else:
                return 'Firefox%20Setup%20{0}.exe'.format(product_version)
        elif 'osx' == os:
            if 'aurora' in alias or 'nightly' in alias:
                return 'firefox-{0}.en-US.mac.dmg'.format(product_version)
            else:
                return 'Firefox%20{0}.dmg'.format(product_version)
        elif os in('linux', 'linux64'):
            if 'aurora' in alias or 'nightly' in alias:
                return 'firefox-{0}.en-US.linux-i686.tar.bz2'.format(product_version)
            else:
                return 'firefox-{0}.tar.bz2'.format(product_version)
        else:
            e = 'Failed, unsupported OS: os = %s, alias = %s, product_version = %s' % \
                (os, alias, product_version)
            raise ValueError(e)

    def verify_redirect_to_correct_product(self, base_url, fx_pkg_name, get_params):
        """Given a set of GET params, this method verifies Bouncer redirects to
        the expected Firefox product.

        :param base_url: The server under test.
        :param fx_pkg_name: The full Firefox package name that is the expected download.
        :param get_params: The GET params to pass to Bouncer.
        """
        response = self.request_with_headers(base_url, params=get_params)
        parsed_url = urlparse(response.url)
        # verify service is up and a 200 OK is returned
        assert requests.codes.ok == response.status_code, \
            'Redirect failed with HTTP status. %s' % \
            self.response_info_failure_message(base_url, get_params, response)
        # verify download location
        assert parsed_url.netloc in self.cdn_netloc_locations, \
            'Failed, redirected to unknown host. %s' % \
            self.response_info_failure_message(base_url, get_params, response)
        # verify Firefox package name and version
        assert fx_pkg_name in response.url, \
            'Failed: Expected product str did not match what was returned %s' % \
            self.response_info_failure_message(base_url, get_params, response)

    def request_with_headers(self, url, params, user_agent=_user_agent_firefox, locale='en-US'):
        """Make a request that includes 'user-agent', 'accept-language', and 'Connection: close' attributes in the
        request header. Redirects will be followed.
        :param url: URL of server to make the request against
        :param params: GET parameters dict
        :param user_agent: Browser useragent string
        :param locale: header accept-language string
        :return: response object
        """
        headers = {'user-agent': user_agent,
                   'accept-language': locale,
                   'Connection': 'close'}
        try:
            response = requests.head(url, headers=headers, verify=False, timeout=15, params=params, allow_redirects=True)
        except requests.RequestException as e:
            request_url = '%s/?%s' % (url, urlencode(params))
            raise AssertionError('Failing URL: %s redirected to %s Error message: %s' % (request_url, response.url, e))
        return response

    def response_info_failure_message(self, url, param, response):
        """Generate a helpful error message that includes the server URL, GET params,
        and header information.
        return 'Failed on %s \nUsing %s.\n %s' % (url,

        :param url: The URL that was under test.
        :param param: The GET params passed to the URL.
        :param response: The response object that was returned by the server under test.
        :return string: An error message with helpful debug information.
        """
        headers = 'Response Headers: '.join(['%s: %s' % (header, value) for header, value in response.headers.items()])
        r_url = 'Response URL: %s' % (response.url)
        return 'Failed on %s Using %s. %s %s.' % (url, param, r_url, headers)

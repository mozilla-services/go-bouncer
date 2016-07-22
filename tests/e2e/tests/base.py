# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.

from urllib import urlencode

import requests
from bs4 import BeautifulSoup


class Base:

    _user_agent_firefox = ('Mozilla/5.0 (Macintosh; Intel Mac OS X 10.7; '
                           'rv:10.0.1) Gecko/20100101 Firefox/10.0.1')

    def _head_request(self, url, user_agent=_user_agent_firefox,
                      locale='en-US', params=None):
        headers = {'user-agent': user_agent,
                   'accept-language': locale,
                   'Connection': 'close'}

        try:
            r = requests.head(url, headers=headers, verify=False, timeout=15,
                              params=params, allow_redirects=False)
        except requests.RequestException as e:
            request_url = self._build_request_url(url, params)

            raise AssertionError('Failing URL: %s.\nError message: %s' % (request_url, e))

        if r.status_code == 302 and r.headers['Location']:
            try:
                request_url = r.headers['Location']
                r = requests.head(request_url, headers=headers, verify=False,
                                  timeout=15, allow_redirects=True)
            except requests.RequestException as e:
                raise AssertionError('Failing URL: %s.\nError message: %s' % (request_url, e))

        return r

    def _parse_response(self, content):
        return BeautifulSoup(content)

    def response_info_failure_message(self, url, param, response):
        return 'Failed on %s \nUsing %s.\n %s' % (url,
                                                  param,
                                                  self.response_info(response))

    def response_info(self, response):
        url = response.url
        return 'Response URL: %s\n Response Headers:\n %s' \
               % (url, self.get_headers(response))

    def get_headers(self, response):
        return '\n'.join(['%s: %s' % (header, value) for header, value in response.headers.items()])

    def _build_request_url(self, url, params):
        if params:
            return '%s/?%s' % (url, urlencode(params))
        else:
            return url

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.

from urllib import urlencode
from urlparse import urlparse

import pytest
import requests

from . import releng_utils as utils
from .test_base import Base


class TestRedirects(Base):

    _locales = utils.get_firefox_locales()
    _os = ('win', 'win64', 'linux', 'linux64', 'osx')
    _winxp_esr_version = utils.get_version_info_for_alias('firefox-esr-latest')
    _winxp_products = [
        'stub',
        'latest',
        'sha1',
        'esr-latest',
        'esr-stub',
        'beta',
        'beta-latest',
        'beta-sha',
        'beta-stub',
        '38.5.1esr',
        '40.0.0esr',
        '58.0.0esr',
        '42.0',
        '43.0.1',
        '49.0b8',
        '49.0b8-ssl',
        '100.0',
        'cats'
    ]

    @pytest.mark.parametrize(('product_alias'), _winxp_products)
    def test_ie6_winxp_useragent_5_1_redirects_to_correct_version(self, base_url, product_alias):
        user_agent_ie6 = ('Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1)')
        param = {
            'product': 'firefox-' + product_alias,
            'lang': 'en-US',
            'os': 'win'
        }
        response = self.request_with_headers(base_url, user_agent=user_agent_ie6, params=param)
        assert self._winxp_esr_version in response.url, param

    @pytest.mark.parametrize(('product_alias'), _winxp_products)
    def test_ie6_winxp_useragent_5_2_redirects_to_correct_version(self, base_url, product_alias):
        user_agent_ie6 = ('Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.2; SV1)')
        param = {
            'product': 'firefox-' + product_alias,
            'lang': 'en-US',
            'os': 'win'
        }
        response = self.request_with_headers(base_url, user_agent=user_agent_ie6, params=param)
        assert self._winxp_esr_version in response.url, param

    def test_that_checks_redirect_using_incorrect_query_values(self, base_url):
        param = {
            'product': 'firefox-47.0.1',
            'lang': 'kitty_language',
            'os': 'stella'
        }
        response = self.request_with_headers(base_url, params=param)

        assert requests.codes.not_found == response.status_code, \
            self.response_info_failure_message(base_url, param, response)

        parsed_url = urlparse(response.url)

        assert 'http' == parsed_url.scheme, 'Failed to redirect to the correct scheme. %s' % \
            self.response_info_failure_message(base_url, param, response)

        assert urlparse(base_url).netloc == parsed_url.netloc, \
            self.response_info_failure_message(base_url, param, response)

        assert urlencode(param) == parsed_url.query, \
            self.response_info_failure_message(base_url, param, response)

    @pytest.mark.parametrize('os', _os)
    def test_verify_locales_redirect_to_the_expected_product(self, base_url, locale, os, version):
        """Verifies the downloaded version of Firefox matches the expected version number
        and filename when Firefox is requested for a specific OS and locale.
        """
        # Ja locale has a macOS-specific locale
        if locale == 'ja' and os == 'osx':
            locale = 'ja-JP-mac'

        get_params = {
            'product': 'firefox-' + version,
            'lang': locale,
            'os': os
        }

        fx_pkg_name = self.get_expected_fx_pkg_str(os, 'latest', version)
        self.verify_redirect_to_correct_product(base_url, fx_pkg_name, get_params)


def pytest_generate_tests(metafunc):
    if 'locale' in metafunc.fixturenames:
        argvalues = []

        def valid_version(version):
            # there are no builds for aurora locales
            return 'a' not in version

        locales = utils.get_firefox_locales()
        _versions = filter(valid_version, utils.get_product_mappings().values())
        for locale, versions in locales.items():
            for version in [v for v in versions if v in _versions]:
                argvalues.append((locale, version))
        metafunc.parametrize(['locale', 'version'], argvalues)

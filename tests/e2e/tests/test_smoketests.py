# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.

import pytest
import requests
from urlparse import urlparse

from base import Base
from utils import RelengHelper


class TestSmokeTests(Base):

    @pytest.mark.smoketest
    def test_verify_releng_aliases_match_what_we_expect(self):
        """Verify the Firefox product list that Releng maintains matches the expected
        list. These product mappings are incorporated into tests, we need to be made
        aware if the keys that releng provide change.

        The json file being verified:
        https://product-details.mozilla.org/1.0/firefox_versions.json
        """
        releng_obj = RelengHelper()
        releng_aliases = releng_obj.fetch_current_fx_product_details()
        bouncer_aliases = releng_obj.releng_to_bouncer_alias_dict
        assert releng_aliases.keys().sort() == bouncer_aliases.keys().sort()

    @pytest.mark.smoketest
    @pytest.mark.parametrize('os', ('win', 'osx', 'linux'))
    def test_verify_firefox_aliases_redirect_to_correct_products(self, base_url, os):
        """Verifies the downloaded version of Firefox matches the expected version number
        and filename when a resource is requested using a go-bouncer alias.

        The test verifies the following aliases: firefox-latest, firefox-esr-latest,
        firefox-nightly-latest, firefox-beta-latest, firefox-beta-latest, firefox-aurora-latest.
        """
        releng_obj = RelengHelper()
        releng_aliases = releng_obj.fetch_current_fx_product_details()
        expected_products = releng_obj.generate_fx_alias_ver_mappings(releng_aliases)
        get_params = {
            'product': 'alias',
            'lang': 'en-US',
            'os': os
        }
        for alias, product_version in expected_products.iteritems():
            fx_pkg_name = self.get_expected_fx_pkg_str(os, alias, product_version)
            # set the GET params that will be sent to bouncer.
            get_params['product'] = alias
            # make the GET request
            response = self._head_request(base_url, params=get_params)
            parsed_url = urlparse(response.url)
            # verify service is up
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

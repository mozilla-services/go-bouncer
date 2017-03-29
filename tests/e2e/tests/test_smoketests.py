# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.

import pytest
from base import Base
import releng_utils as utils


class TestSmokeTests(Base):

    @pytest.mark.smoketest
    def test_verify_releng_aliases_match_what_we_expect(self):
        """Verify the Firefox product list that Releng maintains matches the expected
        list. These product mappings are incorporated into tests, we need to be made
        aware if the keys that releng provide change.

        The json file being verified:
        https://product-details.mozilla.org/1.0/firefox_versions.json
        """
        releng_aliases = utils.fetch_current_fx_product_details()
        bouncer_aliases = utils.releng_to_bouncer_alias_dict
        assert releng_aliases.keys().sort() == bouncer_aliases.keys().sort()

    @pytest.mark.smoketest
    @pytest.mark.parametrize('os', ('win', 'osx', 'linux'))
    def test_verify_firefox_aliases_redirect_to_correct_products(self, base_url, os):
        """Verifies the downloaded version of Firefox matches the expected version number
        and filename when a resource is requested using a go-bouncer alias.

        The test verifies the following aliases: firefox-latest, firefox-esr-latest,
        firefox-nightly-latest, firefox-beta-latest, firefox-beta-latest, firefox-aurora-latest.
        """
        releng_aliases = utils.fetch_current_fx_product_details()
        expected_products = utils.generate_fx_alias_ver_mappings(releng_aliases)
        get_params = {
            'product': 'alias',
            'lang': 'en-US',
            'os': os
        }
        for alias, product_version in expected_products.iteritems():
            fx_pkg_name = self.get_expected_fx_pkg_str(os, alias, product_version)
            # set the GET params that will be sent to bouncer.
            get_params['product'] = alias
            self.verify_redirect_to_correct_product(base_url, fx_pkg_name, get_params)

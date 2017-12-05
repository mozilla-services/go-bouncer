# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.

import pytest

from . import releng_utils as utils
from .test_base import Base


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
    @pytest.mark.xfail(reason="Our Jenkins CI doesn't have access to ship-it")
    def test_verify_product_details_is_in_sync_with_ship_it(self):
        """Verify https://product-details.mozilla.org/1.0/firefox_versions.json and
        https://ship-it.mozilla.org/1.0/firefox_versions.json are in sync.

        Product-details is a mirror of ship-it. Ship-it is behind vpn. Services
        such as download.mozilla.org and mozilla.org pull product information from
        product-details.
        """
        ship_it = utils.fetch_current_fx_product_details(utils._ship_it_url)
        product_details = utils.fetch_current_fx_product_details()
        assert ship_it == product_details, 'ship-it and product-details are out of sync'

    @pytest.mark.smoketest
    def test_verify_firefox_aliases_redirect_to_correct_products(self, base_url, alias, version, os):
        """Verifies the downloaded version of Firefox matches the expected version number
        and filename when a resource is requested using a go-bouncer alias.

        The test verifies the following aliases: firefox-latest, firefox-esr-latest,
        firefox-nightly-latest, firefox-beta-latest, firefox-beta-latest, firefox-aurora-latest.
        """
        get_params = {
            'product': alias,
            'lang': 'en-US',
            'os': os
        }
        fx_pkg_name = self.get_expected_fx_pkg_str(os, alias, version)
        # set the GET params that will be sent to bouncer.
        self.verify_redirect_to_correct_product(base_url, fx_pkg_name, get_params)


def pytest_generate_tests(metafunc):
    if 'alias' in metafunc.fixturenames:
        products = utils.fetch_current_fx_product_details()
        aliases = utils.generate_fx_alias_ver_mappings(products)
        argvalues = []
        for os in ('win', 'osx', 'linux'):
            for alias, version in aliases.items():
                argvalues.append((alias, version, os))
        metafunc.parametrize(['alias', 'version', 'os'], argvalues)

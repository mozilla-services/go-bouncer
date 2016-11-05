# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.

import requests


class RelengHelper:

    _base_url = 'https://product-details.mozilla.org/1.0/%s'
    _firefox_versions_uri = 'firefox_versions.json'
    # mappings adapted off of https://github.com/mozilla-releng/ship-it/blob/master/kickoff/config.py
    releng_to_bouncer_alias_dict = {
        'FIREFOX_AURORA': 'firefox-aurora-latest',
        'FIREFOX_ESR_NEXT': None,  # no checks to run on this product
        'LATEST_FIREFOX_VERSION': 'firefox-latest',
        'FIREFOX_ESR': 'firefox-esr-latest',
        'FIREFOX_NIGHTLY': 'firefox-nightly-latest',
        'LATEST_FIREFOX_OLDER_VERSION': None,  # no checks to run on this product
        'LATEST_FIREFOX_RELEASED_DEVEL_VERSION': 'firefox-beta-latest',
        'LATEST_FIREFOX_DEVEL_VERSION': 'firefox-beta-latest'
    }

    def generate_fx_alias_ver_mappings(self, releng_products, alias_map=releng_to_bouncer_alias_dict):
        """Parses a dictionary that contain releng product/versions and returns
        a new dictionary that maps go-bouncer aliases to expected Firefox product
        versions.

        As the alias_map is walked, aliases with a value
        of type None will not be included in the returned dictionary.

        The key/val mappings that currently represent the system are
        stored in releng_to_bouncer_alias_dict. Currently there are 2 releng products
        that do not have associated go-bouncer aliases and have val = None.

        :arg product_versions: {string:string} the releng object to be walked.
            releng_product/product_version
        :arg alias_map: {string:string} releng to bouncer alias mappings.
            The default, self.releng_to_bouncer_alias_dict is adapted from
            https://github.com/mozilla-releng/ship-it/blob/master/kickoff/config.py
        :returns: {string:string} with product aliases and their  versions numbers"""
        aliases_and_versions = {}
        # create a dict that has {alias: expected_version_num}
        for product, alias in alias_map.iteritems():
            # values of None represent products that don't have aliases
            if alias is not None:
                aliases_and_versions[alias] = releng_products[product]
        return aliases_and_versions

    def fetch_current_fx_product_details(self):
        """Fetches JSON containing key/val pairings of the current releng aliases
        for Firefox and version numbers as known by Mozilla's Release Engineering Team.

        Release Engineering maintains an up-to-date JSON file with the current
        Firefox release values - https://product-details.mozilla.org/1.0/firefox_versions.json.

        :returns: {string:string} dictionary with releng aliases and version numbers"""
        url = self._base_url % self._firefox_versions_uri
        response = requests.get(url)
        releng_products = response.json()

        return releng_products

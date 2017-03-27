# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.

import requests


class FirefoxLocale:

    def __init__(self, locale, version_info):
        self.lang = locale
        self.versions = self.parse_versions(version_info)

    def parse_versions(self, version_info):
        """Parses a dict and returns a list of Firefox versions,
        ignoring Auroa and esr builds.

        :arg version_info: {string:string} A dictionary with k,v of version and
        os.

        :returns: [string] A list of Firefox versions
        """
        versions = []
        for version, os in version_info.items():
            # remove all aurora builds and esr builds
            if 'a' not in version and 'esr' not in version:
                versions.append(version)
        return versions

    def __repr__(self):
        return self.locale


_base_url = 'https://product-details.mozilla.org/1.0/%s'
_firefox_primary_builds_uri = 'firefox_primary_builds.json'
_firefox_versions_uri = 'firefox_versions.json'
# mappings adapted off of https://github.com/mozilla-releng/ship-it/blob/master/kickoff/config.py
releng_to_bouncer_alias_dict = {
    'FIREFOX_AURORA': 'firefox-aurora-latest',
    'FIREFOX_ESR_NEXT': 'firefox-esr-latest',
    'LATEST_FIREFOX_VERSION': 'firefox-latest',
    'FIREFOX_ESR': None,  # no checks to run on this product
    'FIREFOX_NIGHTLY': 'firefox-nightly-latest',
    'LATEST_FIREFOX_OLDER_VERSION': None,  # no checks to run on this product
    'LATEST_FIREFOX_RELEASED_DEVEL_VERSION': 'firefox-beta-latest',
    'LATEST_FIREFOX_DEVEL_VERSION': 'firefox-beta-latest'
}


def generate_fx_alias_ver_mappings(releng_products, alias_map=releng_to_bouncer_alias_dict):
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
        The default, releng_to_bouncer_alias_dict is adapted from
        https://github.com/mozilla-releng/ship-it/blob/master/kickoff/config.py
    :returns: {string:string} with product aliases and their  versions numbers"""
    aliases_and_versions = {}
    # create a dict that has {alias: expected_version_num}
    for product, alias in alias_map.iteritems():
        # values of None represent products that don't have aliases
        if alias is not None:
            aliases_and_versions[alias] = releng_products[product]
    return aliases_and_versions


def fetch_current_fx_product_details():
    """Fetches JSON containing key/val pairings of the current releng aliases
    for Firefox and version numbers as known by Mozilla's Release Engineering Team.

    Release Engineering maintains an up-to-date JSON file with the current
    Firefox release values - https://product-details.mozilla.org/1.0/firefox_versions.json.

    :returns: {string:string} dictionary with releng aliases and version numbers"""
    url = _base_url % _firefox_versions_uri
    response = requests.get(url)
    response.raise_for_status()
    releng_products = response.json()
    return releng_products


def get_firefox_locales():
    """Fetches build versions for each localization of Firefox from Mozilla's Release
    Engineering Team.

    Release Engineering maintains an up-to-date JSON file with the current
    Firefox release values - https://product-details.mozilla.org/1.0/firefox_primary_builds.json.

    :returns list: [FirefoxLocale objects] a list of FirefoxLocale objects.
    """
    locale_objs = []
    url = _base_url % _firefox_primary_builds_uri
    response = requests.get(url)
    response.raise_for_status()
    locale_data =  response.json()
    for locale in locale_data:
        versions = locale_data[locale]
        locale_objs.append(FirefoxLocale(locale, versions))
    return locale_objs

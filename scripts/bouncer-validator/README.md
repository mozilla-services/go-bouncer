# Bouncer Validator

## Description

This script takes two input files — `PRODUCTS_FILE` and `OS_FILE` — and generates Bouncer URLs for each product and operating system combination. It then uses `curl` to request these URLs across a set of predefined `DOMAINS` (e.g., production, staging, dev).

For each request, the script extracts the HTTP `Location` header (i.e., the redirect target) and compares the responses across environments. If any differences are found between the environments (e.g., staging vs. production), they are logged in `domain_mismatches.log`.

This is a diagnostic script designed to verify that staging and development environments behave consistently with production by returning identical redirect URLs.

## Input Files

The script relies on two input JSON files:

- `products.json`: Contains an array of product alias and language combinations.
- `os.json`: Contains a list of supported operating systems.

These files are generated from SQL queries run against the bouncer database (db:bouncer).

## Generating Input Files from SQL

You can generate the input JSON files using the following SQL queries in CloudSQL Studio:

### `products.json`
```sql
SELECT ma.alias, mpl.language
FROM mirror_aliases ma
LEFT JOIN mirror_products mp ON ma.related_product = mp.name
LEFT JOIN mirror_product_langs mpl ON mp.id = mpl.product_id
WHERE mp.active = 1
LIMIT 100;
```

This will generate a list of products with their aliases and languages that are currently active.

### `os.json`
```sql
SELECT name FROM mirror_os LIMIT 100;
```

This will generate a list of supported OS names.

Export the results to `products.json` and `os.json` respectively in JSON format.

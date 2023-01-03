# 📚 ORKL

### The Community Driven Cyber Threat Intelligence Library

---

[![CodeQL](https://github.com/ORKL/library/workflows/CodeQL/badge.svg)](https://github.com/ORKL/library/actions?query=workflow%3ACodeQL)
![GitHub Discussions](https://img.shields.io/github/discussions/ORKL/library)
![Twitter Follow](https://img.shields.io/twitter/follow/orkleu?style=social)
![GitHub Sponsors](https://img.shields.io/github/sponsors/rhaist)

Collection of original report and metadata files that are used by ORKL to
complement the public library of cyber threat intelligence reports. This
repository is used to add reports that are not covered by any of the existing
sources and to provide metadata hints where the PDF metadata or the source
metadata are not sufficient.

In the `corpus` directory the current state of the ORKL library is represented.

The following files can be found in the directory:

* <$REPORT_HASH>.pdf -> Reports not yet in the ORKL CDN (will be imported periodically)
* <$REPORT_HASH>.yaml -> Metadata files that contain a `cdn:` key with the uri of the
    report in the CDN

## ⚙️ Prerequisites

the following requirements need to be installed to use the tool chain:

* Golang (>1.18)
* Exiftool
* Make

on a Debian/Ubuntu system it should be enough to do:

```shell
apt install -y libimage-exiftool-perl golang make
```

**Note:** make sure `$GOPATH/bin` is added to `$PATH`

## ➕ Add Reports

1. Fork this repository and clone the fork to your workspace

2. Install prerequisites and make sure `olm` works:

    ```shell
    make install

    olm -help
    ```

3. Use either the `-report` or the `-recursive` command to add one or multiple
    reports to the corpus. This will copy the PDF(s) as `<HASH>.pdf` and create
    the metadata files as `<HASH>.yaml` in the corpus directory.

4. Work on the metadata files everything that cannot be parsed from the PDF
    metadata needs to researched and added manually.

5. Check and format the .yaml files by running `make` and also make sure the
    parsed title and other report metadata looks clean and you are happy with
    the entries (verify the publication date is correct - do some research).

6. Commit to a patch branch and open a pull-request. If your metadata files and
    TLP level of the reports looks good your PR will be approved after
    verification.

**Note:** please commit both the PDF(s) and the .yaml files. The github action will
upload everything to the CDN. The PDF files will be removed eventually from the
repo by running `olm -janitor`

## 🧹 Fix existing entry in ORKL library

If you want to provided better metadata to an entry that is already in the ORKL
library you need to:

1. Either bootstrap the .yaml file from the api with `olm -hash <$REPORT_HASH>` or `olm -uuid <$REPORT_UUID>` or request a number of entries that need fixing with `olm -work <$NUMBER>`
2. Edit the .yaml file
3. Open a PR to add the file to the main repo

The following video shows an example workflow:

[![asciicast](https://asciinema.org/a/lJ9kpCtcmb52f5zQzazVuOyqV.svg)](https://asciinema.org/a/lJ9kpCtcmb52f5zQzazVuOyqV)

## 📝 Metadata yaml file

the following table is a reference what every .yaml metadata file should contain:

| Key                   | Value                                                      |
|-----------------------|------------------------------------------------------------|
| sha1_hash             | sha1 hash as utf8 string                                   |
| title                 | title of the document as utf8 string                       |
| authors               | authors separated by comma as utf8 string                  |
| pdf_creation_date     | date as ISO 8601                                           |
| pdf_modification_date | date as ISO 8601                                           |
| publication_date      | date as string formatted YYYY-mm-dd                        |
| file_names            | original file names as array                               |
| file_size             | file size in byte as integer                               |
| reference_urls        | urls where the original file can be downloaded as array    |
| language              | language as ISO 639-1 two letter codes                     |
| cdn                   | link to the PDF file in the ORKL cdn (if already imported) |

**Note:** the `cdn` key will be omitted for files that are waiting to be imported

**Example:**

```yaml
sha1_hash: 99411dadc52675d3d86e217564ae8bb7b900754f
title: "Malware Analysis Series (MAS): Article 6"
authors: "Alexandre Borges"
pdf_creation_date: 2022-11-24T13:09:18Z
pdf_modification_date: 2022-11-24T13:09:18Z
publication_date: "2022-11-24"
file_names:
  - mas_6-1.pdf
file_size: 4048456
reference_urls:
  - https://exploitreversing.files.wordpress.com/2022/11/mas_6-1.pdf
language: EN
cdn: https://pub-7cb8ac806c1b4c4383e585c474a24719.r2.dev/99411dadc52675d3d86e217564ae8bb7b900754f.pdf
```

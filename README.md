# Pdf-Service

A standalone service to generate PDF documents from templates.

It uses HTML templates to generate HTML documents and HTML to PDF engines to convert HTML to PDF.

## Templates

All templates are stored in the `templates` directory.

Each template is also a directory, which contains template files and configurations.

On start-up, the service will go through all directories under the `templates` directory and load them as templates if they contain a `templaterc.json` configuration file.

It uses structured [go html templates](https://pkg.go.dev/html/template) and [JSON](https://www.json.org/json-en.html) configuration files.

See examples in [code/templates](https://github.com/hop-/pdf-service/tree/master/templates).

### Template Configuration File

Each template directory should contain a `templaterc.json` configuration file.

The content of the configuration file is the following:

``` json
{
  "name": "{{template-name}}",
  "cover": "{{cover-html-template-file}}",
  "templates": ["{{doc-template-files}}"],
  "style": "{{styles-css}}",
  "header": "{{header-html}}",
  "footer": "{{footer-html}}"
}
```

Some of these fields can be omitted.

## HTML to PDF Engines

There are two engines currently integrated into the service: `wkhtmltopdf` and `chromedp`.

### Wkhtmltopdf

This engine uses `wkhtmltopdf` library to convert HTML to PDF.

This is the official links to the library: [wkhtmltopdf.org](https://wkhtmltopdf.org/), [github](https://github.com/wkhtmltopdf/wkhtmltopdf)

### Chromedb

This engine uses a headless `chrome` session to visualize the generated HTML document, which is then printed as a PDF.

Official links: [go package](https://pkg.go.dev/github.com/chromedp/chromedp), [github](https://github.com/chromedp/chromedp)

## Env Valiables

There are some environment variables that can be used to configure the service.

* `PDF_SERVICE_ROOT` - the root directory where can be found service assets (`configs` and `templates`). Default is the current direcotry
* `PDF_SERVICE_CONCURRENCY` - Max number of cuncurrent workers which can be spawned to generate the pdf docs. Default is `4`
* `REPORT_SERVICE_ENGINE` - HTML to PDF engine which will be used to generate the docs (`chromedp` or `wkhtmltopdf`). Default is `chromedp`
* `PDF_SERVICE_HTTP_ENABLED` - Enable/disable generation of doc over HTTP(S) service. Default is `true`
* `PDF_SERVICE_HTTP_PORT` - HTTP(S) service port. Default is `3000`
* `PDF_SERVICE_HTTPS` - Use TLS layer for HTTP service. Default is `false`
* `PDF_SERVICE_KEY_FILE` - the key file path when HTTPS is enabled
* `PDF_SERVICE_CERT_FILE` - the certificate file path when HTTPS is enabled
* `PDF_SERVICE_KAFKA_ENABLED` - Enable/disable generation over [kafka](https://kafka.apache.org/). Default is `true`
* `PDF_SERVICE_KAFKA_HOST` - Kafka host. Default `kafka:9092`
* `PDF_SERVICE_CREATE_CONSUMER_TOPICS` - Create consumer topics if they are not exist. Default is `false`
* `PDF_SERVICE_KAFKA_REQUESTS_TOPIC` - The name of the requests topic. Default is `PdfRequests`
* `PDF_SERVICE_KAFKA_RESPONSES_TOPIC` - The name of the responses topic. Default is `PdfResponses`

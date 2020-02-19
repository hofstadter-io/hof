{{#if CTX.package-identifier}}
package {{camel CTX.package-identifier}}
{{else if DslContext.module}}
package {{camel file_ddir}}
{{else if (eq DslContext.name CTX.name)}}
package main
{{else}}
package {{camel file_ddir}}
{{/if}}


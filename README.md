# MediaWatch

### Empowering news organizations to fight disinformation

Despite concentrated efforts to combat disinformation and fake news, many countries still exhibit depressingly low trust in the media. Worse, throughout these countries the role of journalism is often brought into question. The entrenched political affiliation of legacy media casts deep doubts on independent and fair fact-checking. On other hand, digital media operate in an unmapped environment, further complicating the issue of misinformation.

MediaWatch aspires to run a pilot project in Greece (108th place in RSF’s 2022 World Press Freedom Index 2022, 32% in media trust in Reuters Institute Digital News Report, 2021) with the further goal to develop a tool that can be used across different countries and media systems.

### Fake News, Mis/Dis–information, Propaganda, all have; *Networks in Common*

MediaWatch is a real-time network analysis platform which continuously monitors online media outlets and identifies flows of information - potentially detecting bad actors and networks of propaganda, with the use of advanced AI algorithms for online content analysis and classification. Therefore, it makes it possible to group articles in clusters by similarity, claims, quotes, entities, topics or categories (and any other combination of custom features) helping journalists, researchers and fact checkers to drill-down information by similar allegations, and rapidly respond to arising issues, reducing the time devoted on non-journalistic tasks.

To corelate passages or claims within articles we use [go-plagiarism](https://github.com/cvcio/go-plagiarism) as our principal algorithm. Though, **we are not interested in plagiarism itself**, we have found, in our long-term feasibility study, that journalists tend to reproduce passages, claims or articles in full (aka copy-paste), as a process in which **an existing narrative is transformed into multiple similar ones, to extend attention to the agenda and frame**, we call this process **"The Chain of Misinformation"**.

## How it works

MediaWatch comprises of multiple micro-servives. In genereal: we collect data for rss feeds using the [feeds](cmd/feeds) micro-service, from twitter using [listen](cmd/listen) and [twitter](cmd/twitter) (which will deprecate*), scrape the data using the [scraper](cmd/scraper) micro-service, enrich the data using [enrich](cmd/enrich) micro-service and using various [AI models](https://huggingface.co/cvcio), and finally create the relationships usgin [compare](cmd/compare) micro-service. All internall communication between micro-services is handled by [kafka](https://kafka.apache.org/) and the [worker](microservice). To store raw articles we use [elasticsearch](https://www.elastic.co/) -different [indices](pkg/es/indeces/) for each language, for the relationships we use [neo4j](https://neo4j.com/), and for any other data we use [mongodb](https://www.mongodb.com/).

![MediaWatch CORE](./assets/MediaWatch%20Core.drawio.png) 
*micro-services architecture*

## Roadmap

In the -not so- near future we plan to introduce multiple new features and micro-services, starting from a unified subscription model for organisations and users to support a cross-organisation fact-checking scheme, where multiple users can share insights, in conjunction with smart-annotations and reports micro-services. Notably:

- Users and Organizations
- Smart Annotations
- Reports
- Hidden Votes
- Important Features Highlighter (claims, quotes, etc.)
- Fully Integrated Lucene Search
- Data Exports
- Open Source Application
- Public API (gRPC, HTTP)
- Documentation and Manual

*Please if you have any suggestions or feature requests reach us at info@cvcio.org, or via github issues.*

## Contributing

If you're new to contributing to Open Source on Github, [this guide](https://opensource.guide/how-to-contribute/) can help you get started. Please check out the contribution guide for more details on how issues and pull requests work. Before contributing be sure to review the [code of conduct](/CODE_OF_CONDUCT.md).

## License

This library is distributed under the MIT license found in the [LICENSE](/LICENSE) file.
# -*- coding: utf-8 -*-
from setuptools import setup

packages = [
    "server",
    "server.config",
    "server.nlp",
    "server.proto.enrich",
    "server.tests",
]

package_data = {"": ["*"]}

install_requires = [
    "gensim>=3.8.1,<4.0.0",
    "grpcio-status>=1.46.0,<2.0.0",
    "grpcio>=1.46.0,<2.0.0",
    "networkx>=2.8,<3.0",
    "nltk>=3.7,<4.0",
    "python-dateutil>=2.8.2,<3.0.0",
    "python-dotenv>=0.20.0,<0.21.0",
    "spacy-lookups-data>=1.0.3,<2.0.0",
    "spacy>=3.3.0,<4.0.0",
    "torch==1.10.2",
    "transformers>=4.18.0,<5.0.0",
]

setup_kwargs = {
    "name": "server",
    "version": "2.2.0",
    "description": "",
    "long_description": None,
    "author": "andefined",
    "author_email": "dimitris.papaevagelou@andefined.com",
    "maintainer": None,
    "maintainer_email": None,
    "url": None,
    "packages": packages,
    "package_data": package_data,
    "install_requires": install_requires,
    "python_requires": ">=3.8,<4.0",
}


setup(**setup_kwargs)

# Language models

Each model is an averaged perceptron trained on the Universal Dependencies
treebanks named below, emitting the universal tag set. Accuracy is measured on
the treebank's own development and test sections through the production
tagger, so it is what a caller of `tag.ForLanguage` gets on text of that kind.

Written by `build.sh`; do not edit by hand.

| Language | Code | Training data | License | Dev | Test | Size |
| --- | --- | --- | --- | --- | --- | --- |
| Catalan | `ca` | UD_Catalan-AnCora | CC BY 4.0 | 98.03% | 97.82% | 1.3 MB |
| Danish | `da` | UD_Danish-DDT | CC BY-SA 4.0 | 94.22% | 94.23% | 0.8 MB |
| German | `de` | UD_German-GSD | CC BY-SA 4.0 | 94.16% | 93.93% | 1.6 MB |
| Spanish | `es` | UD_Spanish-AnCora | CC BY 4.0 | 98.03% | 97.98% | 1.4 MB |
| Finnish | `fi` | UD_Finnish-TDT | CC BY-SA 4.0 | 92.61% | 92.71% | 1.8 MB |
| French | `fr` | UD_French-GSD | CC BY-SA 4.0 | 97.00% | 96.67% | 1.6 MB |
| Croatian | `hr` | UD_Croatian-SET | CC BY-SA 4.0 | 96.52% | 96.28% | 1.3 MB |
| Indonesian | `id` | UD_Indonesian-GSD | CC BY-SA 4.0 | 93.38% | 93.59% | 1.0 MB |
| Italian | `it` | UD_Italian-MarkIT; UD_Italian-PUD; UD_Italian-Valico | CC BY 4.0, CC BY-SA 3.0, CC BY-SA 4.0 | 92.57% | 92.42% | 0.6 MB |
| Korean | `ko` | UD_Korean-GSD | CC BY-SA 4.0 | 90.04% | 90.73% | 1.4 MB |
| Lithuanian | `lt` | UD_Lithuanian-ALKSNIS | CC BY-SA 4.0 | 92.02% | 89.81% | 0.7 MB |
| Latvian | `lv` | UD_Latvian-LVTB | CC BY-SA 4.0 | 96.06% | 95.77% | 2.0 MB |
| Norwegian Bokmål | `nb` | UD_Norwegian-Bokmaal | CC BY-SA 4.0 | 96.34% | 95.96% | 1.3 MB |
| Dutch | `nl` | UD_Dutch-Alpino | CC BY-SA 4.0 | 95.29% | 93.86% | 1.2 MB |
| Portuguese | `pt` | UD_Portuguese-Bosque | CC BY-SA 4.0 | 95.72% | 95.58% | 1.2 MB |
| Romanian | `ro` | UD_Romanian-RRT | CC BY-SA 4.0 | 96.63% | 96.32% | 1.3 MB |
| Russian | `ru` | UD_Russian-GSD; UD_Russian-Taiga | CC BY-SA 4.0 | 96.38% | 95.65% | 1.1 MB |
| Serbian | `sr` | UD_Serbian-SET | CC BY-SA 4.0 | 96.69% | 97.34% | 0.7 MB |
| Swedish | `sv` | UD_Swedish-Talbanken | CC BY-SA 4.0 | 93.68% | 94.86% | 0.7 MB |
| Turkish | `tr` | UD_Turkish-BOUN | CC BY-SA 4.0 | 89.72% | 90.93% | 1.4 MB |

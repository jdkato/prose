# Language models

Each model is an averaged perceptron trained on the Universal Dependencies
treebanks named below. English, the built-in model, emits the Penn Treebank
tag set and is trained with hyphenated spans joined, as prose tokenizes them;
every other model emits the universal tag set. Accuracy is measured on the
first treebank's own development and test sections through the production
tagger, so it is what a caller of `tag.ForLanguage` gets on text of that kind.

Written by `build.sh`; do not edit by hand.

| Language | Code | Training data | License | Dev | Test | Size |
| --- | --- | --- | --- | --- | --- | --- |
| Danish | `da` | UD_Danish-DDT | CC BY-SA 4.0 | 94.22% | 94.23% | 0.8 MB |
| German | `de` | UD_German-GSD; UD_German-HDT | CC BY-SA 4.0 | 94.49% | 94.92% | 4.3 MB |
| English | `en` | UD_English-EWT | CC BY-SA 4.0 | 93.36% | 93.44% | 1.9 MB |
| Spanish | `es` | UD_Spanish-AnCora; UD_Spanish-GSD | CC BY 4.0, CC BY-SA 4.0 | 97.74% | 97.68% | 2.0 MB |
| Finnish | `fi` | UD_Finnish-TDT | CC BY-SA 4.0 | 92.61% | 92.71% | 1.8 MB |
| French | `fr` | UD_French-GSD | CC BY-SA 4.0 | 97.00% | 96.67% | 1.6 MB |
| Italian | `it` | UD_Italian-MarkIT; UD_Italian-PUD; UD_Italian-Valico | CC BY 4.0, CC BY-SA 3.0, CC BY-SA 4.0 | 92.57% | 92.42% | 0.6 MB |
| Norwegian Bokmål | `nb` | UD_Norwegian-Bokmaal | CC BY-SA 4.0 | 96.34% | 95.96% | 1.3 MB |
| Dutch | `nl` | UD_Dutch-Alpino | CC BY-SA 4.0 | 95.29% | 93.86% | 1.2 MB |
| Portuguese | `pt` | UD_Portuguese-Bosque | CC BY-SA 4.0 | 95.72% | 95.58% | 1.2 MB |
| Swedish | `sv` | UD_Swedish-Talbanken | CC BY-SA 4.0 | 93.68% | 94.86% | 0.7 MB |
| Turkish | `tr` | UD_Turkish-BOUN | CC BY-SA 4.0 | 89.72% | 90.93% | 1.4 MB |

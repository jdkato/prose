import json
import pathlib

from nltk.tokenize import TreebankWordTokenizer, sent_tokenize

if __name__ == '__main__':
    t = TreebankWordTokenizer()
    with open(pathlib.PurePath('testdata', 'tokenize.json')) as d:
        data = json.load(d)

    words = []
    sents = []
    for text in data:
        for s in sent_tokenize(text):
            sents.append(s)
            words.append(t.tokenize(s))

    with open(pathlib.PurePath('testdata', 'treebank_words.json'), 'w') as f:
        json.dump(words, f, indent=4)

    with open(pathlib.PurePath('testdata', 'treebank_sents.json'), 'w') as f:
        json.dump(sents, f, indent=4)

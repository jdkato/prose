import os

from newspaper import Article

url = 'http://fox13now.com/2013/12/30/new-year-new-laws-obamacare-pot-guns-and-drones/'
article = Article(url)

article.download()
article.parse()

with open(os.path.join('testdata', 'article.txt'), 'w') as f:
    f.write(article.text)

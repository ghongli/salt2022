#!/usr/local/bin/python

import sys
import nltk

from nltk.tokenize import word_tokenize

nltk.download('punkt')
nltk.download('stopwords')
stopwords = nltk.corpus.stopwords.words('english')

def main():
    for comment in sys.stdin:
        words_in_comment = word_tokenize(comment)
        comment_no_stopwords = [word for word in words_in_comment if word.lower() not in stopwords]
        print( ' '.join(comment_no_stopwords) + '\n')

if __name__ == '__main__':
    main()

package template

import "newspopper/model"

type Parser func(model.Article) ([]byte, error)

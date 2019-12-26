# BinEncoder

[![codecov](https://codecov.io/gh/milQA/binencoder/branch/master/graph/badge.svg)](https://codecov.io/gh/milQA/binencoder)
[![Build Status](https://travis-ci.org/milQA/binencoder.svg?branch=master)](https://travis-ci.org/milQA/binencoder)

## Instal

```
go get github.com/milQA/binencoder
```

## Info

NewEncoder принимает на вход bytes.Buffer и binary.LittleEndian

```go
buf := new(bytes.Buffer)
encoder := binencoder.NewEncoder(buf, binary.LittleEndian)
```

Encode принимает на вход какую-нибудь структуру и длину байтовой записи.
Если необходимо использовать стандартную для типа длину, необходимо задать = 0.

Длину байтовой структуры для поля структуры можно задать тегом:

```go
 `len:"10"`
```

Если тег не задан, то len наследуется от родительского поля, если он был задан, или принимает значение
заданное в функции Encode, или принимается равным 0.

Для того, чтобы игнорировать длину, заданную для родительского поля, либо задать собственную длину,
необходимо явно указать длину в теге.

Если задать тег: 

```go
`len:"-"`
```

поле будет пропущено.

(!) Логика тегов на данный момент некорректно работает с BigEndian.

Типы, которые он может серилизовать функция: bool, uint8, uint16, uint32, int32, uint64, int64, string, slice, struct.
Серилизация происходить последовательно и зависит от структуры типа.

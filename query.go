package gentt

type QueryResult struct {
  Result []EntityId
}

func (res QueryResult) Each(f func(EntityId)) {
  for _, entity := range res.Result {
    f(entity)
  }
}

func (res QueryResult) First() (EntityId, bool) {
  if len(res.Result) == 0 {
    return 0, false
  }
  return res.Result[0], true
}

func (res QueryResult) All() []EntityId {
  return res.Result
}

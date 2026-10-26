package game

type ZoneView struct {
    Id      Identity
    Desc    string
    Exits   []Identity
    Objects []Identity
}

package types

type EpicAccountId string

func (id EpicAccountId) IsValid() bool {
	return id != ""
}

func (id EpicAccountId) String() string {
	return string(id)
}

type ProductUserId string

func (id ProductUserId) IsValid() bool {
	return id != ""
}

func (id ProductUserId) String() string {
	return string(id)
}

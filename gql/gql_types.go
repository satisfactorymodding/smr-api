package gql

import (
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/util"
)

func GenCompInfoToDBCompInfo(gen *generated.CompatibilityInfoInput) *util.CompatibilityInfo {
	if gen == nil {
		return nil
	}
	return &util.CompatibilityInfo{
		Ea:         GenCompToDBComp(gen.Ea),
		Exp:        GenCompToDBComp(gen.Exp),
		Controller: GenControllerCompToDBControllerComp(gen.Controller),
	}
}

func GenCompToDBComp(gen *generated.CompatibilityInput) util.Compatibility {
	r := util.Compatibility{
		State: string(gen.State),
	}
	SetINN(gen.Note, &r.Note)
	return r
}

func GenControllerCompToDBControllerComp(gen *generated.ControllerCompatibilityInput) util.Compatibility {
	r := util.Compatibility{
		State: string(gen.State),
	}
	SetINN(gen.Note, &r.Note)
	return r
}

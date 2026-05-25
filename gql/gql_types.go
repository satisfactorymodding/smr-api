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

func setStateINN(state string, note *string) util.Compatibility {
	r := util.Compatibility{
		State: state,
	}
	SetINN(note, &r.Note)
	return r
}

func GenCompToDBComp(gen *generated.CompatibilityInput) util.Compatibility {
	return setStateINN(string(gen.State), gen.Note)
}

func GenControllerCompToDBControllerComp(gen *generated.ControllerCompatibilityInput) util.Compatibility {
	return setStateINN(string(gen.State), gen.Note)
}

func GenAIDisclosureInfoToDBAIDisclosureInfo(gen *generated.AIUseDisclosureInput) *util.AIUseDisclosureInfo {
	if gen == nil {
		return nil
	}
	r := &util.AIUseDisclosureInfo{
		DisclosureType: string(gen.DisclosureType),
	}
	SetINN(gen.Message, &r.Message)
	return r
}

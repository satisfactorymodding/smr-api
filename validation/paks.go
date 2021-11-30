package validation

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/Vilsol/ue4pak/parser"

	// Import satisfactory-specific types
	_ "github.com/Vilsol/ue4pak/parser/games/satisfactory"
	"github.com/rs/zerolog/log"
)

var classInheritance = map[string]string{
	// FGItemDescriptor Tree
	"FGItemDescriptor":            "",
	"FGBuildDescriptor":           "FGItemDescriptor",
	"FGEquipmentDescriptor":       "FGItemDescriptor",
	"FGItemDescriptorBiomass":     "FGItemDescriptor",
	"FGItemDescriptorNuclearFuel": "FGItemDescriptor",
	// "FGNoneDescriptor":            "FGItemDescriptor",
	"FGResourceDescriptor":           "FGItemDescriptor",
	"FGWildCardDescriptor":           "FGItemDescriptor",
	"FGResourceSinkCreditDescriptor": "FGItemDescriptor",

	"FGBuildingDescriptor": "FGBuildDescriptor",
	"FGVehicleDescriptor":  "FGBuildDescriptor",

	"FGDecorDescriptor": "FGBuildingDescriptor",
	"FGPoleDescriptor":  "FGBuildingDescriptor",

	"FGConsumableDescriptor": "FGEquipmentDescriptor",
	"FGDecorationDescriptor": "FGEquipmentDescriptor",

	"FGResourceDescriptorGeyser": "FGResourceDescriptor",

	// FGBuildable Tree
	"FGBuildable":                "",
	"FGBuildableConveyorBase":    "FGBuildable",
	"FGBuildableDecor":           "FGBuildable",
	"FGBuildableFactory":         "FGBuildable",
	"FGBuildableFactoryBuilding": "FGBuildable",
	"FGBuildableHubTerminal":     "FGBuildable",
	"FGBuildablePole":            "FGBuildable",
	"FGBuildablePowerPole":       "FGBuildable",
	"FGBuildableRailroadBridge":  "FGBuildable",
	"FGBuildableRailroadTrack":   "FGBuildable",
	"FGBuildableBuildableRoad":   "FGBuildable",
	"FGBuildableSpeedSign":       "FGBuildable",
	"FGBuildableStandaloneSign":  "FGBuildable",
	"FGBuildableWire":            "FGBuildable",

	"FGBuildableConveyorBelt": "FGBuildableConveyorBase",
	"FGBuildableConveyorLift": "FGBuildableConveyorBase",

	"FGBuildableConveyorAttachment":    "FGBuildableFactory",
	"FGBuildableDockingStation":        "FGBuildableFactory",
	"FGBuildableGenerator":             "FGBuildableFactory",
	"FGBuildableManufacturer":          "FGBuildableFactory",
	"FGBuildableRadarTower":            "FGBuildableFactory",
	"FGBuildableRailroadSignal":        "FGBuildableFactory",
	"FGBuildableRailroadSwitchControl": "FGBuildableFactory",
	"FGBuildableResourceExtractor":     "FGBuildableFactory",
	"FGBuildableSpaceElevator":         "FGBuildableFactory",
	"FGBuildableStorage":               "FGBuildableFactory",
	"FGBuildableTradingPost":           "FGBuildableFactory",
	"FGBuildableTrainPlatform":         "FGBuildableFactory",
	"FGBuildableWindTurbine":           "FGBuildableFactory",
	"FGBuildableResourceSink":          "FGBuildableFactory",
	"FGBuildablePipeReservoir":         "FGBuildableFactory",
	"FGBuildableResourceSinkShop":      "FGBuildableFactory",
	"FGBuildablePipePart":              "FGBuildableFactory",

	"FGBuildableFloor":      "FGBuildableFactoryBuilding",
	"FGBuildableFoundation": "FGBuildableFactoryBuilding",
	"FGBuildableWalkway":    "FGBuildableFactoryBuilding",
	"FGBuildableWall":       "FGBuildableFactoryBuilding",

	"FGConveyorPoleStackable": "FGBuildablePole",

	"FGBuildableAttachmentMerger":   "FGBuildableConveyorAttachment",
	"FGBuildableAttachmentSplitter": "FGBuildableConveyorAttachment",
	"FGBuildableSplitterSmart":      "FGBuildableAttachmentSplitter",

	"FGBuildableGeneratorFuel":       "FGBuildableGenerator",
	"FGBuildableGeneratorGeoThermal": "FGBuildableGenerator",
	"FGBuildableGeneratorNuclear":    "FGBuildableGeneratorFuel",

	"FGBuildableAutomatedWorkBench": "FGBuildableManufacturer",
	"FGBuildableConverter":          "FGBuildableManufacturer",

	"FGBuildableCentralStorageContainer": "FGBuildableStorage",

	"FGBuildableRailroadStation":    "FGBuildableTrainPlatform",
	"FGBuildableTrainPlatformCargo": "FGBuildableTrainPlatform",
	"FGBuildableTrainPlatformEmpty": "FGBuildableTrainPlatform",

	"FGBuildableRamp":  "FGBuildableFoundation",
	"FGBuildableStair": "FGBuildableFoundation",

	"FGBuildablePoweredWall": "FGBuildableWall",
	"FGBuildableSignWall":    "FGBuildableWall",

	"FGBuildablePipeHyperPart":  "FGBuildablePipePart",
	"FGBuildablePipeHyperStart": "FGBuildablePipeHyperPart",
	"FGPipeHyperStart":          "FGBuildablePipeHyperPart",

	// FGBuildablePipeBase Tree
	"FGBuildablePipeBase":  "",
	"FGBuildablePipeHype":  "FGBuildablePipeBase",
	"FGBuildablePipeline":  "FGBuildablePipeBase",
	"FGBuildablePipeHyper": "FGBuildablePipeBase",

	// FGRecipe Tree
	"FGRecipe":         "",
	"FGResearchRecipe": "FGRecipe",

	// FGBuildCategory Tree
	"FGBuildCategory":    "",
	"FGBuildSubCategory": "FGBuildCategory",

	// FGResourceNode Tree
	"FGResourceNode":       "",
	"FGResourceNodeGeyser": "FGResourceNode",
	"FGResourceDeposit":    "FGResourceNode",

	// FGUnlock Tree
	"FGUnlock":                  "",
	"FGUnlockRecipe":            "FGUnlock",
	"FGUnlockScannableResource": "FGUnlock",

	// FGCharacterBase Tree
	"FGCharacterBase": "",
	"FGCreature":      "FGCharacterBase",
	"FGEnemy":         "FGCreature",

	// FGBuildablePipelineAttachment Tree
	"FGBuildablePipelineAttachment": "",
	"FGBuildablePipelineJunction":   "FGBuildablePipelineAttachment",
	"FGBuildablePipelinePump":       "FGBuildablePipelineAttachment",

	// FGEquipment Tree
	"FGEquipment":           "",
	"FGConsumableEquipment": "FGEquipment",

	// FGVehicle Tree
	"FGVehicle":         "",
	"FGRailroadVehicle": "FGVehicle",
	"FGLocomotive":      "FGRailroadVehicle",
	"FGFreightWagon":    "FGRailroadVehicle",

	// Root
	// "FGFactoryConnectionComponent": "",
	"FGSchematic":            "",
	"FGResearchTree":         "",
	"FGResearchTreeNode":     "",
	"FGItemCategory":         "",
	"FGSchematicCategory":    "",
	"FGInventoryComponent":   "",
	"FGDamageOverTime":       "",
	"FGResourceSinkSettings": "",
	"FGWorkBench":            "",
	"FGEquipmentAttachment":  "",

	// Special
	"BodySetup": "",
}

func AttemptExtractDataFromPak(ctx context.Context, reader parser.PakReader) (data map[string]map[string][]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s\n%s", r, string(debug.Stack()))
		}
	}()

	return ExtractDataFromPak(ctx, reader), nil
}

// TODO Extract Images
func ExtractDataFromPak(ctx context.Context, reader parser.PakReader) map[string]map[string][]interface{} {
	p := parser.NewParser(reader)

	entries := make([]*parser.PakEntrySet, 0)
	p.ProcessPak(ctx, nil, func(s string, entry *parser.PakEntrySet, _ *parser.PakFile) {
		entries = append(entries, entry)
	})

	exportsData := make(map[string]map[string][]interface{})

	for _, entry := range entries {
		exportData := make(map[string][]interface{})
		for _, export := range entry.Exports {
			if export.Export.TemplateIndex.Reference != nil {
				if imp, ok := export.Export.TemplateIndex.Reference.(*parser.FObjectImport); ok {
					cleanName := trim(imp.ClassName)
					treeSize := GetTreeSize(cleanName)
					if treeSize > 0 || strings.HasPrefix(cleanName, "BP_") {
						if _, ok := exportData[cleanName]; !ok {
							exportData[cleanName] = make([]interface{}, 0)
						}

						exportData[cleanName] = append(exportData[cleanName], DecodePropertyFields(ctx, export.Data.Properties))
					} else {
						if _, ok := ignoredClasses[cleanName]; !ok {
							if !strings.HasPrefix(cleanName, "Widget_") &&
								!strings.HasPrefix(cleanName, "ParticleModule") {
								log.Ctx(ctx).Warn().Msgf("Parsing unknown class name: %s", cleanName)

								if _, ok := exportData[cleanName]; !ok {
									exportData[cleanName] = make([]interface{}, 0)
								}

								exportData[cleanName] = append(exportData[cleanName], DecodePropertyFields(ctx, export.Data.Properties))
							}
						}
					}
				}
			}
		}

		if len(exportData) > 0 {
			exportsData[trim(entry.ExportRecord.FileName)] = exportData
		}
	}

	return exportsData
}

func trim(s string) string {
	return strings.Trim(s, "\x00")
}

func GetTreeSize(className string) int {
	size := 0

	for {
		parent, ok := classInheritance[className]

		if !ok {
			break
		}

		size++
		className = parent
	}

	return size
}

func IsA(child string, parent string) bool {
	for {
		if parent == child {
			return true
		}

		var ok bool
		child, ok = classInheritance[child]

		if !ok {
			return false
		}
	}
}

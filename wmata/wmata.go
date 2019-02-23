package wmata

const (
	MetroCenter               = "A01"
	FarragutNorth             = "A02"
	DupontCircle              = "A03"
	WoodleyParkZooAdamsMorgan = "A04"
	ClevelandPark             = "A05"
	VanNessUDC                = "A06"
	TenleytownAU              = "A07"
	FriendshipHeights         = "A08"
	Bethesda                  = "A09"
	MedicalCenter             = "A10"
	GrosvenorStrathmore       = "A11"
	WhiteFlint                = "A12"
	Twinbrook                 = "A13"
	Rockville                 = "A14"
	ShadyGrove                = "A15"
	JudiciarySquare           = "B02"
	UnionStation              = "B03"
	RhodeIslandAveBrentwood   = "B04"
	BrooklandCUA              = "B05"
	Takoma                    = "B07"
	SilverSpring              = "B08"
	ForestGlen                = "B09"
	Wheaton                   = "B10"
	Glenmont                  = "B11"
	NoMaGallaudetU            = "B35"
	Pentagon                  = "C07"
	PentagonCity              = "C08"
	CrystalCity               = "C09"
	NationalAirport           = "C10"
	BraddockRoad              = "C12"
	KingStOldTown             = "C13"
	EisenhowerAve             = "C14"
	Huntington                = "C15"
	MtVernonSq                = "E01"
	Shaw                      = "E02"
	UStCardozo                = "E03"
	ColumbiaHeights           = "E04"
	GeorgiaAvePetworth        = "E05"
	FortTotten                = "E06"
	WestHyattsville           = "E07"
	PrinceGeorgesPlaza        = "E08"
	CollegeParkUMD            = "E09"
	Greenbelt                 = "E10"
	GalleryPlace              = "F01"
	Archives                  = "F02"
	LEnfant                   = "F03"
	Waterfront                = "F04"
	NavyYard                  = "F05"
	Anacostia                 = "F06"
	CongressHeights           = "F07"
	SouthernAve               = "F08"
	NaylorRoad                = "F09"
	Suitland                  = "F10"
	BranchAve                 = "F11"
)

var (
	YellowLine = [...]string{
		FortTotten,
		GeorgiaAvePetworth,
		ColumbiaHeights,
		UStCardozo,
		Shaw,
		MtVernonSq,
		GalleryPlace,
		Archives,
		LEnfant,
		Pentagon,
		PentagonCity,
		CrystalCity,
		NationalAirport,
		BraddockRoad,
		KingStOldTown,
		EisenhowerAve,
		Huntington,
	}
	GreenLine = [...]string{
		Greenbelt,
		CollegeParkUMD,
		PrinceGeorgesPlaza,
		WestHyattsville,
		FortTotten,
		GeorgiaAvePetworth,
		ColumbiaHeights,
		UStCardozo,
		Shaw,
		MtVernonSq,
		GalleryPlace,
		Archives,
		LEnfant,
		Waterfront,
		NavyYard,
		Anacostia,
		CongressHeights,
		SouthernAve,
		NaylorRoad,
		Suitland,
		BranchAve,
	}
	RedLine = [...]string {
		ShadyGrove,
		Rockville,
		Twinbrook,
		WhiteFlint,
		GrosvenorStrathmore,
		MedicalCenter,
		Bethesda,
		FriendshipHeights,
		TenleytownAU,
		VanNessUDC,
		ClevelandPark,
		WoodleyParkZooAdamsMorgan,
		DupontCircle,
		FarragutNorth,
		MetroCenter,
		GalleryPlace,
		JudiciarySquare,
		UnionStation,
		NoMaGallaudetU,
		RhodeIslandAveBrentwood,
		BrooklandCUA,
		FortTotten,
		Takoma,
		SilverSpring,
		ForestGlen,
		Wheaton,
		Glenmont,
	}
)

// GetTrainsByStation retrieves a list of train predictions for the passed station codes
func GetTrainsByStation(stationCodes []string) {

}

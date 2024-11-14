package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/bytedance/gopkg/lang/fastrand"
)

type Decimal1 = int

type weatherStationSource struct {
	name []byte
	avg  Decimal1
}

const MEASUREMENT_DIVERGENCE = 109

func main() {
	for _, station := range SOURCE_STATIONS {
		s := 0
		n := station.name[:len(station.name)-1]
		for _, b := range []byte(n) {
			s += int(b)
		}
	}
	if len(os.Args) < 2 {
		panic("missing parameter: number of records to create (int)")
	}
	count, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		panic("invalid parameter: number of records to create (int)")
	}
	maxRecordLength := 0
	maxLengthAfterName := len("-99.9\n")
	for _, ws := range SOURCE_STATIONS {
		maxRecordLength = max(maxRecordLength, len(ws.name)+maxLengthAfterName)
	}
	maxFileSize := count * int64(maxRecordLength)
	f, err := os.Create("measurements.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.Truncate(maxFileSize)
	data, err := syscall.Mmap(
		int(f.Fd()),
		0,
		int(maxFileSize),
		syscall.PROT_WRITE,
		syscall.MAP_SHARED|syscall.MAP_POPULATE,
	)
	if err != nil {
		panic(err)
	}
	if err := syscall.Madvise(data, syscall.MADV_SEQUENTIAL); err != nil {
		panic(err)
	}
	defer syscall.Munmap(data)

	written := 0
	nextTick := time.Now().Add(time.Second)
	fmt.Println("Generating measurements.txt...")
	for i := range count {
		if i%10000000 == 0 {
			now := time.Now()
			if now.After(nextTick) {
				fmt.Printf("\r%0.2f%%", (float64(i) / (float64(count) / 100)))
				nextTick = now.Add(time.Second)
			}
		}
		station := SOURCE_STATIONS[fastrand.Int()%len(SOURCE_STATIONS)]
		written += copy(data[written:], station.name)
		measurement := station.avg + fastrand.Int()%(MEASUREMENT_DIVERGENCE*2+1) - MEASUREMENT_DIVERGENCE
		measurement = max(-999, min(999, measurement))
		written += writeMeasurement(data[written:], measurement)
	}
	fmt.Println("\r100.00%")
	f.Truncate(int64(written))
}

func writeMeasurement(data []byte, measurement int) int {
	measurementBuffer := [6]byte{}
	bufPos := 0
	if measurement < 0 {
		measurementBuffer[0] = '-'
		bufPos = 1
		measurement = -measurement
	}
	scale := 100
	if measurement < 100 {
		scale = 10
	}
	for scale > 1 {
		measurementBuffer[bufPos] = '0' + byte(measurement/scale)
		measurement %= scale
		scale /= 10
		bufPos++
	}
	measurementBuffer[bufPos] = '.'
	measurementBuffer[bufPos+1] = '0' + byte(measurement)
	measurementBuffer[bufPos+2] = '\n'
	bufPos += 3
	return copy(data, measurementBuffer[:bufPos])
}

var SOURCE_STATIONS = [...]weatherStationSource{
	{[]byte("Abha;"), 180},
	{[]byte("Abidjan;"), 260},
	{[]byte("Abéché;"), 294},
	{[]byte("Accra;"), 264},
	{[]byte("Addis Ababa;"), 160},
	{[]byte("Adelaide;"), 173},
	{[]byte("Aden;"), 291},
	{[]byte("Ahvaz;"), 254},
	{[]byte("Albuquerque;"), 140},
	{[]byte("Alexandra;"), 110},
	{[]byte("Alexandria;"), 200},
	{[]byte("Algiers;"), 182},
	{[]byte("Alice Springs;"), 210},
	{[]byte("Almaty;"), 100},
	{[]byte("Amsterdam;"), 102},
	{[]byte("Anadyr;"), -69},
	{[]byte("Anchorage;"), 28},
	{[]byte("Andorra la Vella;"), 98},
	{[]byte("Ankara;"), 120},
	{[]byte("Antananarivo;"), 179},
	{[]byte("Antsiranana;"), 252},
	{[]byte("Arkhangelsk;"), 13},
	{[]byte("Ashgabat;"), 171},
	{[]byte("Asmara;"), 156},
	{[]byte("Assab;"), 305},
	{[]byte("Astana;"), 35},
	{[]byte("Athens;"), 192},
	{[]byte("Atlanta;"), 170},
	{[]byte("Auckland;"), 152},
	{[]byte("Austin;"), 207},
	{[]byte("Baghdad;"), 228},
	{[]byte("Baguio;"), 195},
	{[]byte("Baku;"), 151},
	{[]byte("Baltimore;"), 131},
	{[]byte("Bamako;"), 278},
	{[]byte("Bangkok;"), 286},
	{[]byte("Bangui;"), 260},
	{[]byte("Banjul;"), 260},
	{[]byte("Barcelona;"), 182},
	{[]byte("Bata;"), 251},
	{[]byte("Batumi;"), 140},
	{[]byte("Beijing;"), 129},
	{[]byte("Beirut;"), 209},
	{[]byte("Belgrade;"), 125},
	{[]byte("Belize City;"), 267},
	{[]byte("Benghazi;"), 199},
	{[]byte("Bergen;"), 77},
	{[]byte("Berlin;"), 103},
	{[]byte("Bilbao;"), 147},
	{[]byte("Birao;"), 265},
	{[]byte("Bishkek;"), 113},
	{[]byte("Bissau;"), 270},
	{[]byte("Blantyre;"), 222},
	{[]byte("Bloemfontein;"), 156},
	{[]byte("Boise;"), 114},
	{[]byte("Bordeaux;"), 142},
	{[]byte("Bosaso;"), 300},
	{[]byte("Boston;"), 109},
	{[]byte("Bouaké;"), 260},
	{[]byte("Bratislava;"), 105},
	{[]byte("Brazzaville;"), 250},
	{[]byte("Bridgetown;"), 270},
	{[]byte("Brisbane;"), 214},
	{[]byte("Brussels;"), 105},
	{[]byte("Bucharest;"), 108},
	{[]byte("Budapest;"), 113},
	{[]byte("Bujumbura;"), 238},
	{[]byte("Bulawayo;"), 189},
	{[]byte("Burnie;"), 131},
	{[]byte("Busan;"), 150},
	{[]byte("Cabo San Lucas;"), 239},
	{[]byte("Cairns;"), 250},
	{[]byte("Cairo;"), 214},
	{[]byte("Calgary;"), 44},
	{[]byte("Canberra;"), 131},
	{[]byte("Cape Town;"), 162},
	{[]byte("Changsha;"), 174},
	{[]byte("Charlotte;"), 161},
	{[]byte("Chiang Mai;"), 258},
	{[]byte("Chicago;"), 98},
	{[]byte("Chihuahua;"), 186},
	{[]byte("Chișinău;"), 102},
	{[]byte("Chittagong;"), 259},
	{[]byte("Chongqing;"), 186},
	{[]byte("Christchurch;"), 122},
	{[]byte("City of San Marino;"), 118},
	{[]byte("Colombo;"), 274},
	{[]byte("Columbus;"), 117},
	{[]byte("Conakry;"), 264},
	{[]byte("Copenhagen;"), 91},
	{[]byte("Cotonou;"), 272},
	{[]byte("Cracow;"), 93},
	{[]byte("Da Lat;"), 179},
	{[]byte("Da Nang;"), 258},
	{[]byte("Dakar;"), 240},
	{[]byte("Dallas;"), 190},
	{[]byte("Damascus;"), 170},
	{[]byte("Dampier;"), 264},
	{[]byte("Dar es Salaam;"), 258},
	{[]byte("Darwin;"), 276},
	{[]byte("Denpasar;"), 237},
	{[]byte("Denver;"), 104},
	{[]byte("Detroit;"), 100},
	{[]byte("Dhaka;"), 259},
	{[]byte("Dikson;"), -111},
	{[]byte("Dili;"), 266},
	{[]byte("Djibouti;"), 299},
	{[]byte("Dodoma;"), 227},
	{[]byte("Dolisie;"), 240},
	{[]byte("Douala;"), 267},
	{[]byte("Dubai;"), 269},
	{[]byte("Dublin;"), 98},
	{[]byte("Dunedin;"), 111},
	{[]byte("Durban;"), 206},
	{[]byte("Dushanbe;"), 147},
	{[]byte("Edinburgh;"), 93},
	{[]byte("Edmonton;"), 42},
	{[]byte("El Paso;"), 181},
	{[]byte("Entebbe;"), 210},
	{[]byte("Erbil;"), 195},
	{[]byte("Erzurum;"), 51},
	{[]byte("Fairbanks;"), -23},
	{[]byte("Fianarantsoa;"), 179},
	{[]byte("Flores, Petén;"), 264},
	{[]byte("Frankfurt;"), 106},
	{[]byte("Fresno;"), 179},
	{[]byte("Fukuoka;"), 170},
	{[]byte("Gabès;"), 195},
	{[]byte("Gaborone;"), 210},
	{[]byte("Gagnoa;"), 260},
	{[]byte("Gangtok;"), 152},
	{[]byte("Garissa;"), 293},
	{[]byte("Garoua;"), 283},
	{[]byte("George Town;"), 279},
	{[]byte("Ghanzi;"), 214},
	{[]byte("Gjoa Haven;"), -144},
	{[]byte("Guadalajara;"), 209},
	{[]byte("Guangzhou;"), 224},
	{[]byte("Guatemala City;"), 204},
	{[]byte("Halifax;"), 75},
	{[]byte("Hamburg;"), 97},
	{[]byte("Hamilton;"), 138},
	{[]byte("Hanga Roa;"), 205},
	{[]byte("Hanoi;"), 236},
	{[]byte("Harare;"), 184},
	{[]byte("Harbin;"), 50},
	{[]byte("Hargeisa;"), 217},
	{[]byte("Hat Yai;"), 270},
	{[]byte("Havana;"), 252},
	{[]byte("Helsinki;"), 59},
	{[]byte("Heraklion;"), 189},
	{[]byte("Hiroshima;"), 163},
	{[]byte("Ho Chi Minh City;"), 274},
	{[]byte("Hobart;"), 127},
	{[]byte("Hong Kong;"), 233},
	{[]byte("Honiara;"), 265},
	{[]byte("Honolulu;"), 254},
	{[]byte("Houston;"), 208},
	{[]byte("Ifrane;"), 114},
	{[]byte("Indianapolis;"), 118},
	{[]byte("Iqaluit;"), -93},
	{[]byte("Irkutsk;"), 10},
	{[]byte("Istanbul;"), 139},
	{[]byte("İzmir;"), 179},
	{[]byte("Jacksonville;"), 203},
	{[]byte("Jakarta;"), 267},
	{[]byte("Jayapura;"), 270},
	{[]byte("Jerusalem;"), 183},
	{[]byte("Johannesburg;"), 155},
	{[]byte("Jos;"), 228},
	{[]byte("Juba;"), 278},
	{[]byte("Kabul;"), 121},
	{[]byte("Kampala;"), 200},
	{[]byte("Kandi;"), 277},
	{[]byte("Kankan;"), 265},
	{[]byte("Kano;"), 264},
	{[]byte("Kansas City;"), 125},
	{[]byte("Karachi;"), 260},
	{[]byte("Karonga;"), 244},
	{[]byte("Kathmandu;"), 183},
	{[]byte("Khartoum;"), 299},
	{[]byte("Kingston;"), 274},
	{[]byte("Kinshasa;"), 253},
	{[]byte("Kolkata;"), 267},
	{[]byte("Kuala Lumpur;"), 273},
	{[]byte("Kumasi;"), 260},
	{[]byte("Kunming;"), 157},
	{[]byte("Kuopio;"), 34},
	{[]byte("Kuwait City;"), 257},
	{[]byte("Kyiv;"), 84},
	{[]byte("Kyoto;"), 158},
	{[]byte("La Ceiba;"), 262},
	{[]byte("La Paz;"), 237},
	{[]byte("Lagos;"), 268},
	{[]byte("Lahore;"), 243},
	{[]byte("Lake Havasu City;"), 237},
	{[]byte("Lake Tekapo;"), 87},
	{[]byte("Las Palmas de Gran Canaria;"), 212},
	{[]byte("Las Vegas;"), 203},
	{[]byte("Launceston;"), 131},
	{[]byte("Lhasa;"), 76},
	{[]byte("Libreville;"), 259},
	{[]byte("Lisbon;"), 175},
	{[]byte("Livingstone;"), 218},
	{[]byte("Ljubljana;"), 109},
	{[]byte("Lodwar;"), 293},
	{[]byte("Lomé;"), 269},
	{[]byte("London;"), 113},
	{[]byte("Los Angeles;"), 186},
	{[]byte("Louisville;"), 139},
	{[]byte("Luanda;"), 258},
	{[]byte("Lubumbashi;"), 208},
	{[]byte("Lusaka;"), 199},
	{[]byte("Luxembourg City;"), 93},
	{[]byte("Lviv;"), 78},
	{[]byte("Lyon;"), 125},
	{[]byte("Madrid;"), 150},
	{[]byte("Mahajanga;"), 263},
	{[]byte("Makassar;"), 267},
	{[]byte("Makurdi;"), 260},
	{[]byte("Malabo;"), 263},
	{[]byte("Malé;"), 280},
	{[]byte("Managua;"), 273},
	{[]byte("Manama;"), 265},
	{[]byte("Mandalay;"), 280},
	{[]byte("Mango;"), 281},
	{[]byte("Manila;"), 284},
	{[]byte("Maputo;"), 228},
	{[]byte("Marrakesh;"), 196},
	{[]byte("Marseille;"), 158},
	{[]byte("Maun;"), 224},
	{[]byte("Medan;"), 265},
	{[]byte("Mek'ele;"), 227},
	{[]byte("Melbourne;"), 151},
	{[]byte("Memphis;"), 172},
	{[]byte("Mexicali;"), 231},
	{[]byte("Mexico City;"), 175},
	{[]byte("Miami;"), 249},
	{[]byte("Milan;"), 130},
	{[]byte("Milwaukee;"), 89},
	{[]byte("Minneapolis;"), 78},
	{[]byte("Minsk;"), 67},
	{[]byte("Mogadishu;"), 271},
	{[]byte("Mombasa;"), 263},
	{[]byte("Monaco;"), 164},
	{[]byte("Moncton;"), 61},
	{[]byte("Monterrey;"), 223},
	{[]byte("Montreal;"), 68},
	{[]byte("Moscow;"), 58},
	{[]byte("Mumbai;"), 271},
	{[]byte("Murmansk;"), 06},
	{[]byte("Muscat;"), 280},
	{[]byte("Mzuzu;"), 177},
	{[]byte("N'Djamena;"), 283},
	{[]byte("Naha;"), 231},
	{[]byte("Nairobi;"), 178},
	{[]byte("Nakhon Ratchasima;"), 273},
	{[]byte("Napier;"), 146},
	{[]byte("Napoli;"), 159},
	{[]byte("Nashville;"), 154},
	{[]byte("Nassau;"), 246},
	{[]byte("Ndola;"), 203},
	{[]byte("New Delhi;"), 250},
	{[]byte("New Orleans;"), 207},
	{[]byte("New York City;"), 129},
	{[]byte("Ngaoundéré;"), 220},
	{[]byte("Niamey;"), 293},
	{[]byte("Nicosia;"), 197},
	{[]byte("Niigata;"), 139},
	{[]byte("Nouadhibou;"), 213},
	{[]byte("Nouakchott;"), 257},
	{[]byte("Novosibirsk;"), 17},
	{[]byte("Nuuk;"), -14},
	{[]byte("Odesa;"), 107},
	{[]byte("Odienné;"), 260},
	{[]byte("Oklahoma City;"), 159},
	{[]byte("Omaha;"), 106},
	{[]byte("Oranjestad;"), 281},
	{[]byte("Oslo;"), 57},
	{[]byte("Ottawa;"), 66},
	{[]byte("Ouagadougou;"), 283},
	{[]byte("Ouahigouya;"), 286},
	{[]byte("Ouarzazate;"), 189},
	{[]byte("Oulu;"), 27},
	{[]byte("Palembang;"), 273},
	{[]byte("Palermo;"), 185},
	{[]byte("Palm Springs;"), 245},
	{[]byte("Palmerston North;"), 132},
	{[]byte("Panama City;"), 280},
	{[]byte("Parakou;"), 268},
	{[]byte("Paris;"), 123},
	{[]byte("Perth;"), 187},
	{[]byte("Petropavlovsk-Kamchatsky;"), 19},
	{[]byte("Philadelphia;"), 132},
	{[]byte("Phnom Penh;"), 283},
	{[]byte("Phoenix;"), 239},
	{[]byte("Pittsburgh;"), 108},
	{[]byte("Podgorica;"), 153},
	{[]byte("Pointe-Noire;"), 261},
	{[]byte("Pontianak;"), 277},
	{[]byte("Port Moresby;"), 269},
	{[]byte("Port Sudan;"), 284},
	{[]byte("Port Vila;"), 243},
	{[]byte("Port-Gentil;"), 260},
	{[]byte("Portland (OR);"), 124},
	{[]byte("Porto;"), 157},
	{[]byte("Prague;"), 84},
	{[]byte("Praia;"), 244},
	{[]byte("Pretoria;"), 182},
	{[]byte("Pyongyang;"), 108},
	{[]byte("Rabat;"), 172},
	{[]byte("Rangpur;"), 244},
	{[]byte("Reggane;"), 283},
	{[]byte("Reykjavík;"), 43},
	{[]byte("Riga;"), 62},
	{[]byte("Riyadh;"), 260},
	{[]byte("Rome;"), 152},
	{[]byte("Roseau;"), 262},
	{[]byte("Rostov-on-Don;"), 99},
	{[]byte("Sacramento;"), 163},
	{[]byte("Saint Petersburg;"), 58},
	{[]byte("Saint-Pierre;"), 57},
	{[]byte("Salt Lake City;"), 116},
	{[]byte("San Antonio;"), 208},
	{[]byte("San Diego;"), 178},
	{[]byte("San Francisco;"), 146},
	{[]byte("San Jose;"), 164},
	{[]byte("San José;"), 226},
	{[]byte("San Juan;"), 272},
	{[]byte("San Salvador;"), 231},
	{[]byte("Sana'a;"), 200},
	{[]byte("Santo Domingo;"), 259},
	{[]byte("Sapporo;"), 89},
	{[]byte("Sarajevo;"), 101},
	{[]byte("Saskatoon;"), 33},
	{[]byte("Seattle;"), 113},
	{[]byte("Ségou;"), 280},
	{[]byte("Seoul;"), 125},
	{[]byte("Seville;"), 192},
	{[]byte("Shanghai;"), 167},
	{[]byte("Singapore;"), 270},
	{[]byte("Skopje;"), 124},
	{[]byte("Sochi;"), 142},
	{[]byte("Sofia;"), 106},
	{[]byte("Sokoto;"), 280},
	{[]byte("Split;"), 161},
	{[]byte("St. John's;"), 50},
	{[]byte("St. Louis;"), 139},
	{[]byte("Stockholm;"), 66},
	{[]byte("Surabaya;"), 271},
	{[]byte("Suva;"), 256},
	{[]byte("Suwałki;"), 72},
	{[]byte("Sydney;"), 177},
	{[]byte("Tabora;"), 230},
	{[]byte("Tabriz;"), 126},
	{[]byte("Taipei;"), 230},
	{[]byte("Tallinn;"), 64},
	{[]byte("Tamale;"), 279},
	{[]byte("Tamanrasset;"), 217},
	{[]byte("Tampa;"), 229},
	{[]byte("Tashkent;"), 148},
	{[]byte("Tauranga;"), 148},
	{[]byte("Tbilisi;"), 129},
	{[]byte("Tegucigalpa;"), 217},
	{[]byte("Tehran;"), 170},
	{[]byte("Tel Aviv;"), 200},
	{[]byte("Thessaloniki;"), 160},
	{[]byte("Thiès;"), 240},
	{[]byte("Tijuana;"), 178},
	{[]byte("Timbuktu;"), 280},
	{[]byte("Tirana;"), 152},
	{[]byte("Toamasina;"), 234},
	{[]byte("Tokyo;"), 154},
	{[]byte("Toliara;"), 241},
	{[]byte("Toluca;"), 124},
	{[]byte("Toronto;"), 94},
	{[]byte("Tripoli;"), 200},
	{[]byte("Tromsø;"), 29},
	{[]byte("Tucson;"), 209},
	{[]byte("Tunis;"), 184},
	{[]byte("Ulaanbaatar;"), -04},
	{[]byte("Upington;"), 204},
	{[]byte("Ürümqi;"), 74},
	{[]byte("Vaduz;"), 101},
	{[]byte("Valencia;"), 183},
	{[]byte("Valletta;"), 188},
	{[]byte("Vancouver;"), 104},
	{[]byte("Veracruz;"), 254},
	{[]byte("Vienna;"), 104},
	{[]byte("Vientiane;"), 259},
	{[]byte("Villahermosa;"), 271},
	{[]byte("Vilnius;"), 60},
	{[]byte("Virginia Beach;"), 158},
	{[]byte("Vladivostok;"), 49},
	{[]byte("Warsaw;"), 85},
	{[]byte("Washington), D.C.;"), 146},
	{[]byte("Wau;"), 278},
	{[]byte("Wellington;"), 129},
	{[]byte("Whitehorse;"), -01},
	{[]byte("Wichita;"), 139},
	{[]byte("Willemstad;"), 280},
	{[]byte("Winnipeg;"), 30},
	{[]byte("Wrocław;"), 96},
	{[]byte("Xi'an;"), 141},
	{[]byte("Yakutsk;"), -88},
	{[]byte("Yangon;"), 275},
	{[]byte("Yaoundé;"), 238},
	{[]byte("Yellowknife;"), -43},
	{[]byte("Yerevan;"), 124},
	{[]byte("Yinchuan;"), 90},
	{[]byte("Zagreb;"), 107},
	{[]byte("Zanzibar City;"), 260},
	{[]byte("Zürich;"), 93},
}

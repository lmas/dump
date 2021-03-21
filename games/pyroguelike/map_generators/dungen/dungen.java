/******************************************************************************
ARGS:
x = kartans x storlek
y = kartans y storlek
objects = antalet objekt att försöka skapa
room_chance = chans att få ett rum
corridor_chance = chans för att få en koriddor

in_file = läs in en fil och visa kartan
out_file =spara genererad karta i en fil

verbose = visa extra text
help = visar hjälp instruktioner
*******************************************************************************
TODO:
1:a
*justera huvudloopen, inte bra med en begränsad foor-loop.. bättre idé?
*lägg till trapporna och andra saker(kistor? etc)
*lägg till argumenten

2:a
*utöka programmet med andra saker man kan göra, istället för att bara generera nya kartor:
	t.ex spara/ladda genererade kartor

3:e
*utöka gränssnittet, GUI?
*olika tilesets, som lava-, is och vattengrottor, öken, skogar?
*undersök andra map format som andra program/spel använder sig av
*******************************************************************************
BUGS:
BUG#001: 
	--CODE--
	newy = getRand(1, ysize-1);
	----
	när newy == 24 crashar programmet "array index out of bounds", något med map variabeln (storleken är definerad som x80, y24)
	--FIX--
	newy = getRand(1, ysize-2); // drar ner newy intervallet till 0<y<24, från 0<y<25
	----
*******************************************************************************
HIGHSCORE:
53/100p på en 80x25 (07-11-30) - utan korridorer
133/200p på en 80x25 (07-11-20) - med korridorer, w00t w00t!!!1 xD måste dra ner på korridorerna dock.. det är dom som spammar upp all yta ¬¬'
95/200p på en 80x25 (07-12-05) - efter uppdatering på korridorerna, chans att få rum: 75%, korridor: 25% (SATANS BRA grafiskt resultat, w00t!!1)
96/200p på en 80x25 (07-12-14) - rum: 75%, korridor: 25%
******************************************************************************/

import java.lang.Integer; //för att vi ska kunna använda Integer.parseInt()
import java.util.*; //för att få dagens "datum"

public class dungen{
	//vad en konsol brukar klara visa som mest
	private int xmax = 80; //80 columns
	private int ymax = 25; //25 rows

	//storleken på kartan
	private int xsize = 0;
	private int ysize = 0;
	
	//antalet "objekt" att generera på kartan
	private int objects = 0;
	
	//bestäm den procentuela chansen att generera antingen ett rum eller en korridor
	//BTW, rum är prioriterade över korridorer så det räcker igentligen med att bara
	//definera rummens "chans"
	private int chanceRoom = 75; 
	private int chanceCorridor = 25;

	//listan som innehåller varje del i kartan
	private int[] dungeon_map = new int[0];
	
	//sparar den gammla seeden för slumpgeneratorn
	private long oldseed = 0;
	
	//definera de olika tiles:en som används
	final private int tileUnused = 0;
	final private int tileDirtWall = 1; //oanvänd, en jord vägg(som man senare ska kunna gräva igenom)
	final private int tileDirtFloor = 2; //jord golv
	final private int tileStoneWall = 3; //opasserbar stenvägg
	final private int tileCorridor = 4;
	final private int tileDoor = 5;
	final private int tileUpStairs = 6;
	final private int tileDownStairs = 7;
	final private int tileChest = 8;
	
	//lite andra, blandade inställningar för programmet nu
	private boolean verbose = false;
	//in_file
	//out_file
	
	//diverse olika meddelanden
	private String msgXSize = "X size of dungeon: \t";
	private String msgYSize = "Y size of dungeon: \t";
	private String msgMaxObjects = "max # of objects: \t";
	private String msgNumObjects = "# of objects made: \t";
	private String msgInFile = "";
	private String msgOutFile = "";
	private String msgHelp = "";
	private String msgDetailedHelp = "";

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	
	//constructor
	public dungen(){
		//note the missing method type (ie., void), not used for constructors
	}

	//funktion för att sätta en kartruta
	private void setCell(int x, int y, int celltype){
		dungeon_map[x + xsize * y] = celltype;
	}
	
	//funktion för att få en kartruta
	private int getCell(int x, int y){
		return dungeon_map[x + xsize * y];
	}
	
	//slumpgeneratorn, vars seed är baserad på sekunderna sedan "java epochen"(tror jag..)
	//antaligen samma datum som unix epochen
	private int getRand(int min, int max){
		Date now = new Date();
		long seed = now.getTime() + oldseed;
		oldseed = seed;
		
		Random randomizer = new Random(seed);
		int n = max - min + 1;
		int i = randomizer.nextInt() % n;
		if (i < 0)
			i = -i;
		//System.out.println("seed: " + seed + "\tnum:  " + (min + i));
		return min + i;
	}
	
	private boolean makeCorridor(int x, int y, int lenght, int direction){
		//TODO: lägga till argument för att bestämma tiletyp som korridoren ska bestå av?
		
		//definerar korridorens dimensioner( erm.. bara längden och riktningen då.. är ju bara en ruta "bredd"..)
		int len = getRand(2, lenght);
		int floor = tileCorridor;
		int dir = 0;
		if (direction > 0 && direction < 4) dir = direction;
		int xtemp = 0;
		int ytemp = 0;
		
		switch(dir){
		case 0:
		//norr
			//kolla om det finns plats för korridoren
			//börja med att kolla att x positionen inte är utanför "map grid:en"
			if (x < 0 || x > xsize) return false;
			else xtemp = x;
			
			//sen kör vi på att kolla om det finns plats för hela korridoren i sin fulla längd och att den inte hamnar utanför "grid:en"
			for (ytemp = y; ytemp > (y-len); ytemp--){
				if (ytemp < 0 || ytemp > ysize) return false;
				if (getCell(xtemp, ytemp) != tileUnused) return false; //gör det inte så skickar vi false och avslutar funktionen negativt
			}
			
			//är vi fortfarande kvar i funktionen, är det fritt fram för korridoren
			for (ytemp = y; ytemp > (y-len); ytemp--){
				setCell(xtemp, ytemp, floor);
			}
			break;
		case 1:
		//öster
				//kolla Ys position, om den är utanför "grid:en" eller inte
				if (y < 0 || y > ysize) return false;
				else ytemp = y;
				
				//sen är det Xs tur
				for (xtemp = x; xtemp < (x+len); xtemp++){
					if (xtemp < 0 || xtemp > xsize) return false;
					if (getCell(xtemp, ytemp) != tileUnused) return false; //gör det inte så skickar vi false och avslutar funktionen
				}
				
				//tydligen gick det bra, så då kan vi lägga ut korridoren
				for (xtemp = x; xtemp < (x+len); xtemp++){
					setCell(xtemp, ytemp, floor);
				}
			break;
		case 2:
		//söder
			//kolla först x positionen att den inte är utanför "grid:en"
			if (x < 0 || x > xsize) return false;
			else xtemp = x;
			
			//sen är det Ys tur..
			for (ytemp = y; ytemp < (y+len); ytemp++){
				if (ytemp < 0 || ytemp > ysize) return false;
				if (getCell(xtemp, ytemp) != tileUnused) return false; //gör det inte så skickar vi false och avslutar funktionen
			}
			
			//lägg ut korridoren sen..
			for (ytemp = y; ytemp < (y+len); ytemp++){
				setCell(xtemp, ytemp, floor);
			}
			break;
		case 3:
		//väster
			if (ytemp < 0 || ytemp > ysize) return false;
			else ytemp = y;
			
			for (xtemp = x; xtemp > (x-len); xtemp--){
				if (xtemp < 0 || xtemp > xsize) return false;
				if (getCell(xtemp, ytemp) != tileUnused) return false; //gör det inte så skickar vi false och avslutar funktionen
			}
		
			for (xtemp = x; xtemp > (x-len); xtemp--){
				setCell(xtemp, ytemp, floor);
			}
			break;
		}
		
		//är vi fortfarande kvar i funktionen ända här nere så tar vi och avslutar den på ett positivt sätt, woot:
		return true;
	}
	
	private boolean makeRoom(int x, int y, int xlength, int ylength, int direction){
		//TODO: gör om till en loop? (eh..? vaddå loop? usch för glömska... ¬¬')
		//TODO: lägg till arguemnt så man kan bestämma tile typerna?
		
		//definera dimensionerna på rummet, ska vara minst 4x4 för rum(ger 2x2 golv att gå på, resten är väggar)
		int xlen = getRand(4, xlength);
		int ylen = getRand(4, ylength);
		//och tile typen som den ska fyllas med
		int floor = tileDirtFloor; //jordgolv..
		int wall = tileDirtWall; //jordvägg
		//bestämmer rummets riktning
		int dir = 0;
		if (direction > 0 && direction < 4) dir = direction;

		switch(dir){
		case 0:
		//norr
			//kolla om det finns plats för det nya rummet
			for (int ytemp = y; ytemp > (y-ylen); ytemp--){
				if (ytemp < 0 || ytemp > ysize) return false;
				for (int xtemp = (x-xlen/2); xtemp < (x+(xlen+1)/2); xtemp++){
					if (xtemp < 0 || xtemp > xsize) return false;
					if (getCell(xtemp, ytemp) != tileUnused) return false; //gör det inte så skickar vi false och avslutar funktionen
				}
			}
			
			//är vi fortfarande kvar i funktionen, skapa rummet
			for (int ytemp = y; ytemp > (y-ylen); ytemp--){
				for (int xtemp = (x-xlen/2); xtemp < (x+(xlen+1)/2); xtemp++){
					//sätt upp rummets väggar först
					if (xtemp == (x-xlen/2)) setCell(xtemp, ytemp, wall);
					else if (xtemp == (x+(xlen-1)/2)) setCell(xtemp, ytemp, wall);
					else if (ytemp == y) setCell(xtemp, ytemp, wall);
					else if (ytemp == (y-ylen+1)) setCell(xtemp, ytemp, wall);
					//sen fyller resten med floor
					else setCell(xtemp, ytemp, floor);
				}
			}
			break;
		case 1:
		//öster
			//kolla om det finns plats för det nya rummet
			for (int ytemp = (y-ylen/2); ytemp < (y+(ylen+1)/2); ytemp++){
				if (ytemp < 0 || ytemp > ysize) return false;
				for (int xtemp = x; xtemp < (x+xlen); xtemp++){
					if (xtemp < 0 || xtemp > xsize) return false;
					if (getCell(xtemp, ytemp) != tileUnused) return false; //gör det inte så skickar vi false och avslutar funktionen
				}
			}

			//är vi fortfarande kvar i funktionen, skapa rummet
			for (int ytemp = (y-ylen/2); ytemp < (y+(ylen+1)/2); ytemp++){
				for (int xtemp = x; xtemp < (x+xlen); xtemp++){
					//sätt upp rummets väggar först
					if (xtemp == x) setCell(xtemp, ytemp, wall);
					else if (xtemp == (x+xlen-1)) setCell(xtemp, ytemp, wall);
					else if (ytemp == (y-ylen/2)) setCell(xtemp, ytemp, wall);
					else if (ytemp == (y+(ylen-1)/2)) setCell(xtemp, ytemp, wall);
					//sen fyller resten med floor
					else setCell(xtemp, ytemp, floor);
				}
			}
			break;
		case 2:
		//söder
			//kolla om det finns plats för det nya rummet
			for (int ytemp = y; ytemp < (y+ylen); ytemp++){
				if (ytemp < 0 || ytemp > ysize) return false;
				for (int xtemp = (x-xlen/2); xtemp < (x+(xlen+1)/2); xtemp++){
					if (xtemp < 0 || xtemp > xsize) return false;
					if (getCell(xtemp, ytemp) != tileUnused) return false; //gör det inte så skickar vi false och avslutar funktionen
				}
			}
			
			//är vi fortfarande kvar i funktionen, skapa rummet
			for (int ytemp = y; ytemp < (y+ylen); ytemp++){
				for (int xtemp = (x-xlen/2); xtemp < (x+(xlen+1)/2); xtemp++){
					//sätt upp rummets väggar först
					if (xtemp == (x-xlen/2)) setCell(xtemp, ytemp, wall);
					else if (xtemp == (x+(xlen-1)/2)) setCell(xtemp, ytemp, wall);
					else if (ytemp == y) setCell(xtemp, ytemp, wall);
					else if (ytemp == (y+ylen-1)) setCell(xtemp, ytemp, wall);
					//sen fyller resten med floor
					else setCell(xtemp, ytemp, floor);
				}
			}
			break;
		case 3:
		//väster
			//kolla om det finns plats för det nya rummet
			for (int ytemp = (y-ylen/2); ytemp < (y+(ylen+1)/2); ytemp++){
				if (ytemp < 0 || ytemp > ysize) return false;
				for (int xtemp = x; xtemp > (x-xlen); xtemp--){
					if (xtemp < 0 || xtemp > xsize) return false;
					if (getCell(xtemp, ytemp) != tileUnused) return false; //gör det inte så skickar vi false och avslutar funktionen
				}
			}
			
			//är vi fortfarande kvar i funktionen, skapa rummet
			for (int ytemp = (y-ylen/2); ytemp < (y+(ylen+1)/2); ytemp++){
				for (int xtemp = x; xtemp > (x-xlen); xtemp--){
					//sätt upp rummets väggar först
					if (xtemp == x) setCell(xtemp, ytemp, wall);
					else if (xtemp == (x-xlen+1)) setCell(xtemp, ytemp, wall);
					else if (ytemp == (y-ylen/2)) setCell(xtemp, ytemp, wall);
					else if (ytemp == (y+(ylen-1)/2)) setCell(xtemp, ytemp, wall);
					//sen fyller resten med floor
					else setCell(xtemp, ytemp, floor);
				}
			}
			break;
		}
		
		//allt gick som det skulle..
		return true;
	}
	
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	
	//funktion för att visa kartan på skärmen
	public void showDungeon(){
		for (int y = 0; y < ysize; y++){
			for (int x = 0; x < xsize; x++){
				//System.out.print(getCell(x, y));
				switch(getCell(x, y)){
				case tileUnused:
					System.out.print(" ");
					break;
				case tileDirtWall:
					System.out.print("+");
					break;
				case tileDirtFloor:
					System.out.print(".");
					break;
				case tileStoneWall:
					System.out.print("O");
					break;
				case tileCorridor:
					System.out.print("#");
					break;
				case tileDoor:
					System.out.print("D");
					break;
				case tileUpStairs:
					System.out.print("<");
					break;
				case tileDownStairs:
					System.out.print(">");
					break;
				case tileChest:
					System.out.print("*");
					break;
				};
			}
			if (xsize < xmax) System.out.print("\n");
		}
	}
	
	//och en funktion för att generera en ny karta
	public boolean createDungeon(int inx, int iny, int inobj){
		if (inobj < 1) objects = 10;
		else objects = inobj;
		
		//justera kartans storlek, om den är större eller mindre än "gränserna"
		if (inx < 3) xsize = 3;
		else if (inx > xmax) xsize = xmax;
		else xsize = inx;
		
		if (iny < 3) ysize = 3;
		else if (iny > ymax) ysize = ymax;
		else ysize = iny;
		
		System.out.println(msgXSize + xsize);
		System.out.println(msgYSize + ysize);
		System.out.println(msgMaxObjects + objects);
		
		//redeklarera kart variabel så att vi får rätt minnesstorlek på den
		dungeon_map = new int[xsize * ysize];
		
		//börja med att sätta upp en standard karta
		for (int y = 0; y < ysize; y++){
			for (int x = 0; x < xsize; x++){
				//dessa fyra rader sätter upp en "mur" runt hela området
				if (y == 0) setCell(x, y, tileStoneWall);
				else if (y == ysize-1) setCell(x, y, tileStoneWall);
				else if (x == 0) setCell(x, y, tileStoneWall);
				else if (x == xsize-1) setCell(x, y, tileStoneWall);
				
				//och fyller allt inom kanterna med "jord"
				else setCell(x, y, tileUnused);
			}
		}
		
		/*******************************************************************************
		Här börjar huvudalgoritmen för generatorn
		*******************************************************************************/
		
		//börja med att göra ett rum nånstans i mitten av kartan, som vi sedan kan utgå från
		//(TODO: istället för ett rum, slumpa fram objekt?)
		makeRoom(xsize/2, ysize/2, 8, 6,getRand(0,3)); //getrand saken för att slumpa fram riktning på rummet
		
		//håller reda på hur många "features" som finns
		int currentFeatures = 1; //+1 för det första rummet i mitten vi skapade först
		
		//sen börjar vi huvud loopen
		for (int countingTries = 0; countingTries < 1000; countingTries++){
			//kolla om vi har skapat max antalet av "object", avsluta om vi har det
			if (currentFeatures == objects){
				break;
			}
				
			//börja med att välja en slumpmässig vägg
			int newx = 0;
			int xmod = 0;
			int newy = 0;
			int ymod = 0;
			int validTile = -1;
			//1000 chanser att hitta en passande plats för den nya saken, heh...
			//går rätt snabbt ändån i java
			for (int testing = 0; testing < 1000; testing++){
				newx = getRand(1, xsize-1);
				newy = getRand(1, ysize-1);
				validTile = -1;
				//System.out.println("tempx: " + newx + "\ttempy: " + newy);
				if (getCell(newx, newy) == tileDirtWall || getCell(newx, newy) == tileCorridor){
					//kolla om vi kan nå denna position
					if (getCell(newx, newy+1) == tileDirtFloor || getCell(newx, newy+1) == tileCorridor){
						validTile = 0; //
						xmod = 0;
						ymod = -1;
					}
					else if (getCell(newx-1, newy) == tileDirtFloor || getCell(newx-1, newy) == tileCorridor){
						validTile = 1; //
						xmod = +1;
						ymod = 0;
					}
					else if (getCell(newx, newy-1) == tileDirtFloor || getCell(newx, newy-1) == tileCorridor){
						validTile = 2; //
						xmod = 0;
						ymod = +1;
					}
					else if (getCell(newx+1, newy) == tileDirtFloor || getCell(newx+1, newy) == tileCorridor){
						validTile = 3; //
						xmod = -1;
						ymod = 0;
					}
					
					//kolla att det inte finns någon dörr i närheten, så vi slipper få en massa korridorer 
					//alldeles bredvid varandra
					if (validTile > -1){
						if (getCell(newx, newy+1) == tileDoor) //norr
							validTile = -1;
						else if (getCell(newx-1, newy) == tileDoor)//öster
							validTile = -1;
						else if (getCell(newx, newy-1) == tileDoor)//söder
							validTile = -1;
						else if (getCell(newx+1, newy) == tileDoor)//väster
							validTile = -1;
					}
					
					//kan vi, hoppar vi ut ur loopen och fortsätter med att skapa saken
					if (validTile > -1) break;
				}
			}
			if (validTile > -1){
				//välj vad som ska skapas nu, och åt vilken riktning
				int feature = getRand(0, 100);
				if (feature <= chanceRoom){ //ett nytt rum
					if (makeRoom((newx+xmod), (newy+ymod), 8, 6, validTile)){
						currentFeatures++;
						
						//sen markerars punkten på väggen med en dörr
						setCell(newx, newy, tileDoor);
						
						//sen måste vi "göra rent" framför dörren så man kan nå den
						setCell((newx+xmod), (newy+ymod), tileDirtFloor);
					}
				}
				else if (feature >= chanceRoom){ //en ny korridor
					if (makeCorridor((newx+xmod), (newy+ymod), 6, validTile)){
						//samma sak som för rum: plussa på summan och fixa två dörrar
						currentFeatures++;
						
						setCell(newx, newy, tileDoor);
					}
				}
			}
		}
		
		
		/*******************************************************************************
		körningen av algoritmen gick bra, avsluta det hela nu
		*******************************************************************************/
		
		//sprid ut "bonussakerna"(trappor, kistor osv.) över kartan nu
		int newx = 0;
		int newy = 0;
		int ways = 0; //från hur många håll vi kan nå den slumpade positionen
		int state = 0; //vilket läge loopen befinner sig i, börja med trapporna
		while (state != 10){
			for (int testing = 0; testing < 1000; testing++){
				//if (state == 10) break;
				newx = getRand(1, xsize-1);
				//--BUG#001--
				newy = getRand(1, ysize-2); // drar ner newy intervallet till 0<y<24, från 0<y<25
				//----
				System.out.println("x: " + newx + "\ty: " + newy);
				ways = 4; //antalet håll man kan nå saken från, ju lägre desto fler håll kan man nå den från
				
				//kolla den nya slumpade positionen om den är nårbar
				if (getCell(newx, newy+1) == tileDirtFloor || getCell(newx, newy+1) == tileCorridor){
				//norr
					if (getCell(newx, newy+1) != tileDoor)
					ways--;
				}
				if (getCell(newx-1, newy) == tileDirtFloor || getCell(newx-1, newy) == tileCorridor){
				//öster
					if (getCell(newx-1, newy) != tileDoor)
					ways--;
				}
				if (getCell(newx, newy-1) == tileDirtFloor || getCell(newx, newy-1) == tileCorridor){
				//söder
					if (getCell(newx, newy-1) != tileDoor)
					ways--;
				}
				if (getCell(newx+1, newy) == tileDirtFloor || getCell(newx+1, newy) == tileCorridor){
				//väster
					if (getCell(newx+1, newy) != tileDoor)
					ways--;
				}
				
				if (state == 0){
					if (ways == 0){
					//placera ut en trappa upp
						setCell(newx, newy, tileUpStairs);
						System.out.println("STATE:\t" + state);
						System.out.println("WAYS:\t" + ways);
						state = 1;
						break;
					}
				}
				else if (state == 1){
					if (ways == 0){
					//placera en trappa ner
						setCell(newx, newy, tileDownStairs);
						System.out.println("STATE:\t" + state);
						System.out.println("WAYS:\t" + ways);
						state = 10;
						break;
					}
				}
			}
		}
		
		
		//allt gick bra, tala om det för användaren
		System.out.println(msgNumObjects + currentFeatures);
		
		//Avslutar funktionen NU med ett positivt resultat
		return true;
	}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	
	public static void main(String[] args){
		//kollar om användaren vill själv bestämma storleken på kartan och
		//antalet "grejer" som ska genereras, t.ex rum och korridorer
		int x = 80; int y = 25; int dungeon_objects = 0;
		
		//konvertera string argumentet till en int, om det finns ett/flera argument
		if (args.length >= 1)
			dungeon_objects = Integer.parseInt(args[0]);

		if (args.length >= 2)
			x = Integer.parseInt(args[1]);
		
		if (args.length >= 3)
			y = Integer.parseInt(args[2]);

		//skapa en ny klass av typ "dungen", och tala om hur många grejer vi vill generera
		dungen generator = new dungen();
		
		//efter att ha initierat classen, så börjar vi med att skapa en ny slumpmässigt genererad dungeon med denna metod
		if (generator.createDungeon(x, y, dungeon_objects)){
			//sen kan det vara bra att visa resultate på skärmen osså...
			generator.showDungeon();
		}
	}
}

/******************************************************************************
by Mike ( mike@mikera.net)

The algorithm
=============
In this algorithm a "feature" is taken to mean any kind of map component
e.g. large room, small room, corridor, circular arena, vault etc.

1.  Fill the whole map with solid earth
2.  Dig out a single room in the centre of the map
3.  Pick a wall of any room
4.  Decide upon a new feature to build
5.  See if there is room to add the new feature through the chosen wall
6.  If yes, continue. If no, go back to step 3
7.  Add the feature through the chosen wall
8.  Go back to step 3, until the dungeon is complete
9.  Add the up and down staircases at random points in map
10. Finally, sprinkle some monsters and items liberally over dungeon

Step 1 and 2 are easy once you've got the map set up. I have found it very
useful to write a "fillRect" command that fills a rectangular map area
with a specified tile type. 

Step 3 is trickier. You can't pick random squares to add new features
because the rule is to always add to the existing dungeon. This makes the
connections look good, and also guarantees that every square is reachable.
The way Tyrant does it is to pick random squares on the map until it finds
a wall square adjacent (horizontally or vertically) to a clear square.
This is a good method, since it gives you a roughly even chance of picking
any particular wall square.

Step 4 isn't too hard - I just use a random switch statement to determine
which of a range of features to build. You can change the whole look of
the map by weighting the probabilities of different features. Well-ordered
dungeons will have lots of regular rooms and long straight corridors. Cave
complexes will tend to have jagged caverns, twisting passages etc.

Step 5 is more tricky, and the key to the whole algorithm. For each
feature, you need to know the area of the map that it will occupy. Then
you need to scan outwards from the chosen wall to see if this area
intersects any features that are already there. Tyrant does this in a
fairly simplistic way - it just works out the rectangular space that the
new feature will occupy plus a square on each side for walls, then checks
to see if the entire rectangle is currently filled with solid earth.

Step 6 decides whether or not to add the feature. If the area under
consideration contains anything other than solid earth already, then the
routine loops back to step 3. Note that *most* new features will be
rejected in this way. This isn't a problem, as the processing time is
negligible. Tyrant tries to add 300 or so features to each dungeon, but
usually only 40 or so make it past this stage.

Step 7 draws the new feature once you've decided the area is OK. In this
stage, you can also add any interesting room features, such as
inhabitants, traps, secret doors and treasure.

Step 8 just loops back to build more rooms. The exact number of times that
you want to do this will depend on map size and various other factors. 

Step 9 is pretty self-explanatory. Easiest way to do it is to write a
routine that picks random squares until it finds an empty one where the
staircases can be added.

Step 10 just creates a host of extra random monsters in random locations
to add some spice. Tyrant creates about most of the monsters at this point
of the map generation, although it does add a few special creatures when
individual rooms are generated.
******************************************************************************/
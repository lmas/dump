//Fractal map generation

public class fracmap{
	//storleken måste vara "a power of 2"
	private int size = 64;
	private int[] fractal_map = new int[size * size];
	
	//för att kunna visa debuggnigns text
	private boolean debug = false;

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	
	private int getCell(int x, int y){
		return fractal_map[x + size * y];
	}
	
	private void setCell(int x, int y, int tile){
		fractal_map[x + size * y] = tile;
	}	

	private double getRand(){
		double temp = Math.random();
		//System.out.println(temp);
		return temp;
	}
	
	public void showMap(){
		for (int y = 0; y < size; y++){
			for (int x = 0; x < size; x++){
				//System.out.print(getCell(x, y));
				switch(getCell(x, y)){
				case 0:
					System.out.print("+");
					break;
				case 1:
					System.out.print(" ");
					break;
				};
			}
			if (size < 80) System.out.print("\n");
		}
	}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	
	public boolean makeMap(int x, int y){
		//xsize = x; ysize = y;
		
		setCell(0, 0, 0);
		setCell(size/2, 0, 0);
		setCell(0, size/2, 0);
		setCell(size/2, size/2, 1);
		
		if (debug){
			showMap();
			System.out.println("------------------------------------------");
		}
		
		makeFractal(size/4);
		
		//System.out.println("done making map");
		return true;
	}

	public boolean makeFractal(int details){
		for (int tempy = 0; tempy < size; tempy += details){
			for (int tempx = 0; tempx < size; tempx += details){
				
				int cx = tempx + ((getRand() < 0.5) ? 0 : details);
				int cy = tempy + ((getRand() < 0.5) ? 0 : details);
				
				cx = (cx / (details*2)) * details * 2;
				cy = (cy / (details*2)) * details * 2;
				
				cx = cx % size;
				cy = cy % size;
				
				setCell(tempx, tempy, getCell(cx, cy));
			}
		}
		if (details > 1){
			if (debug){
				showMap();
				System.out.println("------------------------------------------");
			}
			makeFractal(details/2);
		}
		
		return true;
	}
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	
	public static void main(String[] args){
		fracmap generator = new fracmap();

		//int x = 80; int y = 25; int details = 2;
		int x = 64; int y = 16;
		
		if (generator.makeMap(x, y)){
			generator.showMap();
		}
	}
}

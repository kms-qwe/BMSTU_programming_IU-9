package letuchka1;

public class Main {
    public static void main(String[] args) {
        int[] array = {1, 2, 3, 4, 5, 6, 7, 8, 9};

        NamedMatrix namedMatrix = new NamedMatrix("Matrix A", "Symmetric", array);
        
        namedMatrix.print();

        int sum = namedMatrix.sumDiagonal();
        System.out.println("Sum of diagonal elements: " + sum);
        namedMatrix.resetMatrix();
        System.out.println("Matrix is reset.");
        namedMatrix.print();
    }
}

// 4 варинат
class Matrix {
    public String matrixType;
    public int[][] matrix; 

    public Matrix(String matrixType, int[] array) {
        this.matrixType = matrixType;
        this.matrix = new int[3][3];
        for (int i = 0; i < 3; i++) {
            for (int j = 0; j < 3; j++) {
                this.matrix[i][j] = array[i * 3 + j];
            }
       }
    }
    public void resetMatrix() {
        for (int i = 0; i < 3; i++) {
            for (int j =0 ; j< 3; j++) {
                this.matrix[i][j] = 0;
            }
        }
    }
    public void print() {
        System.out.println("Unnamed matrix");
        for (int i = 0;i<3;i++){
            for (int j = 0; j <3;j++){
                System.out.print(this.matrix[i][j] + " ");
            }
            System.out.println();
        }
    }
}
class NamedMatrix extends Matrix {
    public String matrixName;

    public NamedMatrix(String matrixName, String matrixType, int[] array) {
        super(matrixType, array);
        this.matrixName = matrixName;
    }
    public int sumDiagonal() {
        int sum = 0;
        for (int i = 0; i < 3; i++) {
            sum += this.matrix[i][i];
        }
        return sum; 
    }
     public void print() {
        System.out.println("Name:" +  this.matrixName);
        for (int i = 0;i<3;i++){
            for (int j = 0; j <3;j++){
                System.out.print(this.matrix[i][j] + " ");
            }
            System.out.println();
        }
    }
}
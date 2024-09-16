package lab3;
public class MatrixWithRangSort implements Comparable<MatrixWithRangSort> {
    private int[][] matrix;
    private int m;
    private int n;
    private int rang;

    public MatrixWithRangSort(int[] array, int m, int n) {
        if (array.length != m * n) {
            System.err.println("Error: Array length must be equal to m * n");
            System.exit(1);
        }
        this.m = m;
        this.n = n;
        this.matrix = new int[m][n];
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                this.matrix[i][j] = array[i * n + j];
            }
        }
        this.rang = calculateRank();
    }

    private int calculateRank() {
        int rank = 0;
        int rowCount = m;
        int colCount = n;
        int[][] copyMatrix = new int[rowCount][colCount];

        for (int i = 0; i < rowCount; i++) {
            System.arraycopy(matrix[i], 0, copyMatrix[i], 0, colCount);
        }

        for (int col = 0; col < colCount; col++) {
            int leadingEntryRow = -1;
            for (int i = col; i < rowCount; i++) {
                if (copyMatrix[i][col] != 0) {
                    leadingEntryRow = i;
                    break;
                }
            }

            if (leadingEntryRow != -1) {
                if (leadingEntryRow != col) {
                    int[] temp = copyMatrix[col];
                    copyMatrix[col] = copyMatrix[leadingEntryRow];
                    copyMatrix[leadingEntryRow] = temp;
                }

                for (int i = col + 1; i < rowCount; i++) {
                    int factor = copyMatrix[i][col] / copyMatrix[col][col];
                    for (int j = col; j < colCount; j++) {
                        copyMatrix[i][j] -= factor * copyMatrix[col][j];
                    }
                }
            }
        }

        for (int i = 0; i < rowCount; i++) {
            boolean isZeroRow = true;
            for (int j = 0; j < colCount; j++) {
                if (copyMatrix[i][j] != 0) {
                    isZeroRow = false;
                    break;
                }
            }
            if (!isZeroRow) {
                rank++;
            }
        }

        return rank;
    }

    public int compareTo(MatrixWithRangSort other) {
        return rang - other.rang;
    }

    public void printMatrix() {
        System.out.println("Matrix dimensions: " + m + "x" + n);
        System.out.println("Rank of the matrix: " + rang);
        System.out.println("Matrix:");
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                System.out.print(matrix[i][j] + " ");
            }
            System.out.println();
        }
    }

}


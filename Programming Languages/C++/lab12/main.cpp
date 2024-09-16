#include <iostream>
#include <fstream>
#include <string>
#include <vector>
#include <filesystem>

void numberHeadings(std::string& fileContent) {
    int headingLevel = 0;
    int headingCount = 0;
    std::string::size_type pos = 0;
    while ((pos = fileContent.find("\n", pos)) != std::string::npos) {
        std::string::size_type startPos = pos + 1;
        std::string::size_type endPos = fileContent.find("\n", startPos);
        if (endPos == std::string::npos) {
            endPos = fileContent.length();
        }
        std::string line = fileContent.substr(startPos, endPos - startPos);
        if (line.find("=") != std::string::npos ) {
            headingLevel = 1;
            headingCount++;
            fileContent.replace(startPos, endPos - startPos, line + "\n" + std::to_string(headingCount) + ". ");
        } else if (line.find("-") != std::string::npos) {
            headingLevel = 2;
            headingCount++;
            fileContent.replace(startPos, endPos - startPos, line + "\n" + std::to_string(headingCount) + ". ");
        }
        pos = endPos;
    }
}

int main(int argc, char* argv[]) {
    if (argc < 2) {
        std::cerr << "Usage: " << argv[0] << " <directory_path>" << std::endl;
        return 1;
    }

    std::string directory = argv[1];
    for (const auto& entry : std::filesystem::directory_iterator(directory)) {
        if (entry.path().extension() == ".md") {
            std::ifstream inputFile(entry.path());
            if (!inputFile.is_open()) {
                std::cerr << "Could not open the file: " << entry.path() << std::endl;
                continue;
            }

            std::string fileContent((std::istreambuf_iterator<char>(inputFile)),
                                    std::istreambuf_iterator<char>());
            inputFile.close();

            numberHeadings(fileContent);

            std::ofstream outputFile(entry.path());
            if (!outputFile.is_open()) {
                std::cerr << "Could not write to the file: " << entry.path() << std::endl;
                continue;
            }
            outputFile << fileContent;
            outputFile.close();
        }
    }
    return 0;
}

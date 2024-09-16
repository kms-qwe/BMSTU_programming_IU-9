#include <iostream>
#include <fstream>
#include <string>
#include <filesystem>

void removeScripts(std::string& fileContent) {
    std::string::size_type scriptStart, scriptEnd;

    while ((scriptStart = fileContent.find("<script")) != std::string::npos) {
        scriptEnd = fileContent.find("</script>", scriptStart);
        if (scriptEnd != std::string::npos) {
            scriptEnd += 9; 
            fileContent.erase(scriptStart, scriptEnd - scriptStart);
        } else {
            fileContent.erase(scriptStart, std::string::npos);
        }
    }
}

int main(int argc, char* argv[]) {
    if (argc < 2) {
        std::cerr << "Usage: " << argv[0] << " <directory_path>" << std::endl;
        return 1;
    }

    std::string directory = argv[1];
    for (const auto& entry : std::filesystem::directory_iterator(directory)) {
        if (entry.path().extension() == ".html") {
            std::ifstream inputFile(entry.path());
            if (!inputFile.is_open()) {
                std::cerr << "Could not open the file: " << entry.path() << std::endl;
                continue;
            }

            std::string fileContent((std::istreambuf_iterator<char>(inputFile)),
                                    std::istreambuf_iterator<char>());
            inputFile.close();

            removeScripts(fileContent);

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

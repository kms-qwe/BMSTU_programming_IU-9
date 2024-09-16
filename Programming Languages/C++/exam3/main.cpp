#include <iostream>
#include <fstream>
#include <string>
#include <set>
#include <filesystem>
#include <regex>
#include <algorithm>

namespace fs = std::filesystem;

std::set<std::string> extractTagsFromFile(const std::string& filePath) {
    std::set<std::string> tags;
    std::ifstream file(filePath);
    if (!file.is_open()) {
        std::cerr << "Could not open the file: " << filePath << std::endl;
        return tags;
    }

    std::string line;
    std::regex tagRegex(R"(<\s*([a-zA-Z0-9]+))");
    std::smatch match;

    while (std::getline(file, line)) {
        std::sregex_iterator begin(line.begin(), line.end(), tagRegex);
        std::sregex_iterator end;
        for (std::sregex_iterator i = begin; i != end; ++i) {
            match = *i;
            tags.insert(match[1].str());
        }
    }
    
    file.close();
    return tags;
}

std::set<std::string> findAllHtmlTags(const std::string& directory) {
    std::set<std::string> allTags;

    for (const auto& entry : fs::recursive_directory_iterator(directory)) {
        if (entry.is_regular_file() && entry.path().extension() == ".html") {
            std::set<std::string> fileTags = extractTagsFromFile(entry.path().string());
            allTags.insert(fileTags.begin(), fileTags.end());
        }
    }

    return allTags;
}

void saveTagsToFile(const std::set<std::string>& tags, const std::string& outputFilePath) {
    std::ofstream outFile(outputFilePath);
    if (!outFile.is_open()) {
        std::cerr << "Could not open the file: " << outputFilePath << std::endl;
        return;
    }

    for (const auto& tag : tags) {
        outFile << tag << std::endl;
    }

    outFile.close();
}

int main(int argc, char* argv[]) {
    if (argc < 2) {
        std::cerr << "Usage: " << argv[0] << " <directory_path>" << std::endl;
        return 1;
    }

    std::string directoryPath = argv[1];

    std::set<std::string> allTags = findAllHtmlTags(directoryPath);
    saveTagsToFile(allTags, "tags.txt");

    std::cout << "Tags have been successfully saved to tags.txt" << std::endl;

    return 0;
}

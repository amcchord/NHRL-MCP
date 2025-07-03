# NHRL MCP Server v1.9.0 Release Notes

## 🚀 Major Feature Release

### 🆕 New Features

#### Wiki Tool Integration
- **Complete NHRL Wiki Access**: New `nhrl_wiki` tool provides full access to the NHRL wiki at https://wiki.nhrl.io
- **Search Functionality**: Search wiki pages by keywords with configurable result limits
- **Page Retrieval**: Get full wiki page content with markup support
- **Page Extracts**: Retrieve plain text summaries of wiki pages
- **Read-Only Safety**: Wiki tool is available in all tool modes (reporting, full-safe, full)

#### Enhanced Stats API
- **Improved Data Handling**: Better statistics processing and data management
- **Performance Optimizations**: Enhanced API response times and reliability
- **Extended Coverage**: More comprehensive stats data integration

#### Build System Improvements
- **Complete Build Overhaul**: Brand new build and notarization script with 430+ lines of improvements
- **Enhanced Notarization**: Improved macOS binary signing and notarization process
- **Better Error Handling**: More robust build process with comprehensive error reporting
- **Automated Testing**: New test scripts for validating functionality

### 🛠️ Technical Improvements

- **Code Quality**: Significant codebase improvements with better error handling
- **Documentation**: Comprehensive documentation for all new features
- **Testing**: New test scripts including wiki tool validation
- **Security**: Enhanced credential management and build security

### 📖 Use Cases

The new wiki tool enables:
- **Rules and Regulations**: Quick access to competition rules and requirements
- **Weight Classes**: Detailed information about robot weight class specifications
- **Safety Requirements**: Complete safety checklists and requirements
- **Tournament Information**: Tournament formats, procedures, and qualification details
- **Technical Specifications**: Detailed technical requirements and specifications

### 📦 Downloads

All binaries are code-signed and notarized for macOS.

| Platform | Architecture | File |
|----------|-------------|------|
| Linux | AMD64 | `nhrl-mcp-server-v1.9.0-linux-amd64.tar.gz` |
| Linux | ARM64 | `nhrl-mcp-server-v1.9.0-linux-arm64.tar.gz` |
| macOS | Intel | `nhrl-mcp-server-v1.9.0-darwin-amd64.tar.gz` |
| macOS | Apple Silicon | `nhrl-mcp-server-v1.9.0-darwin-arm64.tar.gz` |
| Windows | AMD64 | `nhrl-mcp-server-v1.9.0-windows-amd64.zip` |
| Windows | ARM64 | `nhrl-mcp-server-v1.9.0-windows-arm64.zip` |

### 🔐 SHA256 Checksums
```
201c736875710d7a384153bed28b5614d6b44f1d634d54c68bd3784a5aa41f9b  nhrl-mcp-server-v1.9.0-linux-amd64.tar.gz
062dba3439a0cc0fd46b17c2365898e5f160e296326a0737473e72284f46050c  nhrl-mcp-server-v1.9.0-darwin-amd64.tar.gz
5c963b10be2c6000113f8a61a8220044475240b5d3215cee18475d35b9dd8f9a  nhrl-mcp-server-v1.9.0-windows-arm64.zip
51eb71a29f338cc57a0941e9455ff14638924d9e82be5e8d9f3720e0f5836ac2  nhrl-mcp-server-v1.9.0-darwin-arm64.tar.gz
bc6511c90f1147a2ee9407a346da116deb38eb8beaeb0217575949942ad30abe  nhrl-mcp-server-v1.9.0-linux-arm64.tar.gz
23c2cfc578af7b07d1cbd3fccdeba61fe806856ad9c73fe93867a7f9da32b22f  nhrl-mcp-server-v1.9.0-windows-amd64.zip
```

### 🔧 Installation & Usage

1. Download the appropriate binary for your platform
2. Extract the archive
3. Run the MCP server with your preferred configuration
4. Access the new wiki tool through your MCP-compatible client

### 🧪 Testing

Test the new wiki functionality:
```bash
./scripts/test-wiki-tool.sh
```

### 📚 Documentation

- **Wiki Tool**: See `NHRL_WIKI_TOOL_DOCUMENTATION.md` for complete usage guide
- **Stats API**: Review `NHRL_STATS_API_UPDATE_SUMMARY.md` for API changes
- **Build Process**: Updated build scripts with comprehensive documentation

### 🔄 Compatibility

- **Full Backward Compatibility**: All existing tools and features remain unchanged
- **No Breaking Changes**: Safe to upgrade from any previous version
- **Enhanced Functionality**: New features complement existing capabilities

### 💡 What's Next

- More wiki integration features
- Enhanced search capabilities
- Additional data sources
- Performance improvements

---

**Full Changelog**: https://github.com/amcchord/NHRL-MCP/compare/v1.8.1...v1.9.0 
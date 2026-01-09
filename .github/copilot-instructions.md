When creating or editing Go files, follow these strict requirements:

1. **Package Declaration and Import Order (CRITICAL)**
   - The `package` declaration MUST be the first non-comment line in every Go file
   - The import block MUST be AFTER the package declaration
   - Order is ALWAYS: comments/license header → package declaration → imports → code
   - NEVER place imports before the package declaration - this will cause compilation errors
   - Example correct structure:
     ```
     // Copyright header (optional)
     
     package main
     
     import (
         "fmt"
         "strings"
     )
     
     func main() { ... }
     ```

2. **Avoid Duplicate Headers**
   - Each Go file must have exactly ONE package declaration
   - Each Go file must have at most ONE import block
   - When editing files, read the existing structure first to avoid duplicating package/import declarations
   - Use replace_string_in_file or multi_replace_string_in_file to edit existing files, not create_file

3. **File Creation Method**
   - ALWAYS use create_file or replace_string_in_file tools for Go files
   - NEVER use heredocs or echo commands in terminals to write Go files
   - NEVER use `cat > file.go << 'EOF'` patterns for Go code
   - Heredocs often fail due to special characters, backticks, and formatting issues
   - If you must use terminal commands, use the file tools instead

4. **Before Creating Go Files**
   - Check if the file already exists using read_file or file_search
   - If editing, read the current content first to understand the structure
   - Preserve existing package declarations and imports when editing
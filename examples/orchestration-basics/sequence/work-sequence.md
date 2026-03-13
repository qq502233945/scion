# Work Sequence

1. Use the web search tool to find a reliable, free public API that returns random trivia or quotes in JSON format, then write a Python script named `fetch_data.py` that queries this API and saves the raw JSON output exactly as received into a new file called `raw_data.json`.

2. Execute the `fetch_data.py` script using the shell to generate the initial `raw_data.json` file, and then create a second Python script named `process_data.py` that parses this JSON file, extracts the specific string containing the trivia or quote, and appends it as a new line to a text file named `processed_log.txt`.

3. Run `fetch_data.py` followed by `process_data.py` via the shell three separate times to accumulate a few distinct entries in `processed_log.txt`, and subsequently use a shell command like `cat` to output the contents of `processed_log.txt` to verify the data extraction and appending worked correctly.

4. Search the web for best practices on dynamically generating minimal HTML5 documents directly from Python without external dependencies, then write a script named `generate_report.py` that reads all the lines from `processed_log.txt` and writes them into a beautifully formatted, standalone `report.html` file as a stylized unordered list.

5. Execute `generate_report.py` via the shell to produce the final `report.html`, and conclude by using shell commands to bundle `fetch_data.py`, `process_data.py`, `generate_report.py`, and `report.html` into a compressed archive named `project_delivery.tar.gz`.

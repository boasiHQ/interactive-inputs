const fs = require('node:fs');
const path = require('node:path');

const ACTION_DOCS_MARKER = '<!-- action-docs-inputs source="action.yml" -->';

/**
 * Removes one surrounding quote pair from an action metadata scalar value.
 *
 * @param {string} value The raw scalar value from action.yml.
 * @returns {string} The scalar without matching single or double quotes.
 */
function unquote(value) {
  const trimmed = value.trim();

  if (
    (trimmed.startsWith('"') && trimmed.endsWith('"')) ||
    (trimmed.startsWith("'") && trimmed.endsWith("'"))
  ) {
    return trimmed.slice(1, -1);
  }

  return trimmed;
}

/**
 * Escapes Markdown table cell content while preserving inline code formatting.
 *
 * @param {string} value The value to render inside a Markdown table cell.
 * @returns {string} Escaped table cell content.
 */
function escapeCell(value) {
  return value.replace(/`/g, '\\`').replace(/\|/g, '\\|');
}

/**
 * Converts a YAML literal block into the compact default format historically
 * used by the README input table.
 *
 * @param {string[]} lines The literal block lines captured from action.yml.
 * @returns {string} A single-line representation of the block.
 */
function formatLiteralDefault(lines) {
  return lines
    .map((line) => line.replace(/^ {6}/, ''))
    .join(' ')
    .trimEnd();
}

/**
 * Parses the inputs section from this action's metadata file.
 *
 * The parser is intentionally narrow: it supports the action.yml shape used in
 * this repository without introducing another YAML dependency into the toolchain.
 *
 * @param {string} metadata The contents of action.yml.
 * @returns {Array<{name: string, description: string, required: string, defaultValue: string}>} Parsed input rows.
 */
function parseInputs(metadata) {
  const lines = metadata.split(/\r?\n/);
  const inputs = [];
  let inInputs = false;
  let current = null;

  const finishCurrent = () => {
    if (!current) {
      return;
    }

    inputs.push({
      name: current.name,
      description: current.description ?? '',
      required: current.required ?? 'false',
      defaultValue: current.defaultValue ?? '""',
    });
  };

  for (let index = 0; index < lines.length; index += 1) {
    const line = lines[index];

    if (line === 'inputs:') {
      inInputs = true;
      continue;
    }

    if (inInputs && /^[a-z]+:/.test(line)) {
      finishCurrent();
      break;
    }

    if (!inInputs || line.trim() === '') {
      continue;
    }

    const inputMatch = line.match(/^ {2}([A-Za-z0-9_-]+):\s*$/);
    if (inputMatch) {
      finishCurrent();
      current = { name: inputMatch[1] };
      continue;
    }

    if (!current) {
      continue;
    }

    const propertyMatch = line.match(/^ {4}([A-Za-z0-9_-]+):(?:\s*(.*))?$/);
    if (!propertyMatch) {
      continue;
    }

    const [, key, rawValue = ''] = propertyMatch;
    if (key === 'description') {
      current.description = unquote(rawValue);
    }

    if (key === 'required') {
      current.required = unquote(rawValue);
    }

    if (key === 'default') {
      if (rawValue.trim() === '|') {
        const block = [];
        while (index + 1 < lines.length && /^ {6}/.test(lines[index + 1])) {
          index += 1;
          block.push(lines[index]);
        }
        current.defaultValue = formatLiteralDefault(block);
      } else {
        current.defaultValue = unquote(rawValue);
      }
    }
  }

  return inputs;
}

/**
 * Renders the README inputs table between the existing action-docs markers.
 *
 * @param {Array<{name: string, description: string, required: string, defaultValue: string}>} inputs Parsed action inputs.
 * @returns {string} Markdown content for the generated section.
 */
function renderInputsSection(inputs) {
  const rows = inputs.map(
    (input) =>
      `| \`${escapeCell(input.name)}\` | <p>${escapeCell(input.description)}</p> | \`${escapeCell(input.required)}\` | \`${escapeCell(input.defaultValue)}\` |`,
  );

  return [
    ACTION_DOCS_MARKER,
    '## Inputs',
    '',
    '| name | description | required | default |',
    '| --- | --- | --- | --- |',
    ...rows,
    ACTION_DOCS_MARKER,
  ].join('\n');
}

/**
 * Replaces the generated inputs section in README.md with data from action.yml.
 */
function updateReadme() {
  const root = path.resolve(__dirname, '..');
  const actionPath = path.join(root, 'action.yml');
  const readmePath = path.join(root, 'README.md');

  const actionMetadata = fs.readFileSync(actionPath, 'utf8');
  const readme = fs.readFileSync(readmePath, 'utf8');
  const generatedSection = renderInputsSection(parseInputs(actionMetadata));
  const markerStart = readme.indexOf(ACTION_DOCS_MARKER);
  const markerEnd = readme.indexOf(
    ACTION_DOCS_MARKER,
    markerStart + ACTION_DOCS_MARKER.length,
  );

  if (markerStart === -1 || markerEnd === -1) {
    throw new Error('README.md must contain two action-docs input markers');
  }

  const updated = [
    readme.slice(0, markerStart),
    generatedSection,
    readme.slice(markerEnd + ACTION_DOCS_MARKER.length),
  ].join('');

  fs.writeFileSync(readmePath, updated);
}

updateReadme();

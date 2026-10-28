const fs = require('node:fs');
const YAML = require('yaml');

// Keep the existing action-docs markers so prose outside generated sections is retained.
const sections = ['inputs'];

function cell(value) {
  return String(value ?? '')
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('|', '&#124;')
    .replaceAll('`', '&#96;')
    .replace(/\r?\n/g, '<br>');
}

function table(entries, isInput) {
  const rows = Object.entries(entries ?? {});
  if (!rows.length) return '';
  const heading = isInput ? 'Inputs' : 'Outputs';
  const columns = isInput
    ? '| name | description | required | default |\n| --- | --- | --- | --- |'
    : '| name | description |\n| --- | --- |';
  const body = rows.map(([name, value]) => {
    const fields = [`\`${cell(name)}\``, `<p>${cell(value.description)}</p>`];
    if (isInput) {
      if (value.required !== undefined && typeof value.required !== 'boolean') {
        throw new Error(`Input ${name}: required must be a YAML boolean`);
      }
      fields.push(
        `\`${value.required ?? false}\``,
        `\`${cell(value.default === '' ? '""' : (value.default ?? '""'))}\``,
      );
    }
    return `| ${fields.join(' | ')} |`;
  });
  return `## ${heading}\n\n${columns}\n${body.join('\n')}`;
}

function generate(readme, metadata) {
  if (!metadata.name || !metadata.description || !metadata.runs?.using) {
    throw new Error(
      'Action metadata must include name, description and runs.using',
    );
  }
  const content = {
    header: `# ${cell(metadata.name)}`,
    description: `## Description\n\n${cell(metadata.description)}`,
    runs: `## Runs\n\nThis action is a \`${cell(metadata.runs.using)}\` action.`,
    inputs: table(metadata.inputs, true),
    outputs: table(metadata.outputs, false),
  };
  for (const section of sections) {
    const marker = `<!-- action-docs-${section} source="action.yml" -->`;
    const parts = readme.split(marker);
    if (parts.length !== 3) {
      throw new Error(`README must contain exactly two ${section} markers`);
    }
    readme = `${parts[0]}${marker}\n${content[section]}\n${marker}${parts[2]}`;
  }
  return readme;
}

function main(args) {
  for (const arg of args) {
    if (!['--check', '--update'].includes(arg)) {
      throw new Error(`Unknown argument: ${arg}`);
    }
  }
  const metadata = YAML.parse(fs.readFileSync('action.yml', 'utf8'));
  const before = fs.readFileSync('README.md', 'utf8');
  const after = generate(before, metadata);
  if (args.includes('--check')) {
    if (before !== after) {
      throw new Error(
        'README action metadata is stale; run yarn readme --update',
      );
    }
  } else if (before !== after) {
    fs.writeFileSync('README.md', after);
  }
}

module.exports = { generate };

if (require.main === module) {
  try {
    main(process.argv.slice(2));
  } catch (error) {
    console.error(error.message);
    process.exitCode = 1;
  }
}

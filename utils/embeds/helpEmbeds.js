/* global client */

module.exports.reactions = [ 'ℹ', '🖥', '🎵', '🛠', '🚨', '🔥' ];

module.exports[ 'ℹ' ] = function () {
  return {
    'title': '**HJÄLP SIDA**',
    'description': '__Tryck knapparna längst ned för att byta sida__.' +
      '\nDu kan få mer information om ett kommando genom att köra `$help <command>`.\n' +
      '\n:information_source: **--** Denna sida' +
      '\n:desktop: **--** Digital Ungdom kommandon' +
      '\n:musical_note:  **--** Musik kommandon' +
      '\n:tools: **--** Misc kommandon' +
      '\n:rotating_light: **--** Admin kommandon' +
      '\n:fire:  **--** Stäng hjälp sida',
    'color': 4086462
  };
};

module.exports[ '🖥' ] = function () {
  const commands = {};
  let longestCommand = 8,
    longestDescription = 10;

  for ( const command of client.commands.values() ) {
    if ( command.group === 'digitalungdom' ) commands[ command.name ] = command.description;
  }

  for ( let command of Object.keys( commands ) ) {
    if ( command.length > longestCommand ) longestCommand = command.length;
    if ( commands[ command ].length > longestDescription ) longestDescription = commands[ command ].length;
  }

  let description = '```\n';

  description += 'KOMMANDO'.padEnd( longestCommand, ' ' ) + ' | ' + 'FÖRKLARING'.padEnd( longestDescription, ' ' ) + '\n';
  description += '-'.repeat( longestCommand + longestDescription ) + '\n';

  for ( let command of Object.keys( commands ).sort() ) {
    description += command.padEnd( longestCommand, ' ' ) + ' | ' + commands[ command ].padEnd( longestDescription, ' ' ) + '\n';
  }

  description += '```\n';

  return {
    'title': '**DIGITAL UNGDOM**',
    'description': description,
    'color': 4086462
  };
};

module.exports[ '🎵' ] = function () {
  const commands = {};
  let longestCommand = 8,
    longestDescription = 10;

  for ( const command of client.commands.values() ) {
    if ( command.group === 'music' ) commands[ command.name ] = command.description;
  }

  for ( let command of Object.keys( commands ) ) {
    if ( command.length > longestCommand ) longestCommand = command.length;
    if ( commands[ command ].length > longestDescription ) longestDescription = commands[ command ].length;
  }

  let description = '```\n';

  description += 'KOMMANDO'.padEnd( longestCommand, ' ' ) + ' | ' + 'FÖRKLARING'.padEnd( longestDescription, ' ' ) + '\n';
  description += '-'.repeat( longestCommand + longestDescription + 3 ) + '\n';

  for ( let command of Object.keys( commands ).sort() ) {
    description += command.padEnd( longestCommand, ' ' ) + ' | ' + commands[ command ].padEnd( longestDescription, ' ' ) + '\n';
  }

  description += '```\n';

  return {
    'title': '**MUSIK**',
    'description': description,
    'color': 4086462
  };
};

module.exports[ '🛠' ] = function () {
  const commands = {};
  let longestCommand = 8,
    longestDescription = 10;

  for ( const command of client.commands.values() ) {
    if ( command.group === 'misc' ) commands[ command.name ] = command.description;
  }

  for ( let command of Object.keys( commands ) ) {
    if ( command.length > longestCommand ) longestCommand = command.length;
    if ( commands[ command ].length > longestDescription ) longestDescription = commands[ command ].length;
  }

  let description = '```\n';

  description += 'KOMMANDO'.padEnd( longestCommand, ' ' ) + ' | ' + 'FÖRKLARING'.padEnd( longestDescription, ' ' ) + '\n';
  description += '-'.repeat( longestCommand + longestDescription + 3 ) + '\n';

  for ( let command of Object.keys( commands ).sort() ) {
    description += command.padEnd( longestCommand, ' ' ) + ' | ' + commands[ command ].padEnd( longestDescription, ' ' ) + '\n';
  }

  description += '```\n';

  return {
    'title': '**MISC**',
    'description': description,
    'color': 4086462
  };
};

module.exports[ '🚨' ] = function () {
  const commands = {};
  let longestCommand = 8,
    longestDescription = 10;

  for ( const command of client.commands.values() ) {
    if ( command.group === 'admin' ) commands[ command.name ] = command.description;
  }

  for ( let command of Object.keys( commands ) ) {
    if ( command.length > longestCommand ) longestCommand = command.length;
    if ( commands[ command ].length > longestDescription ) longestDescription = commands[ command ].length;
  }

  let description = '```\n';

  description += 'KOMMANDO'.padEnd( longestCommand, ' ' ) + ' | ' + 'FÖRKLARING'.padEnd( longestDescription, ' ' ) + '\n';
  description += '-'.repeat( longestCommand + longestDescription + 3 ) + '\n';

  for ( let command of Object.keys( commands ).sort() ) {
    description += command.padEnd( longestCommand, ' ' ) + ' | ' + commands[ command ].padEnd( longestDescription, ' ' ) + '\n';
  }

  description += '```\n';

  return {
    'title': '**ADMIN**',
    'description': description,
    'color': 4086462
  };
};
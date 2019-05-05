/* global client guild memberProcess config */

module.exports = async function ( message ) {
  if ( message.author.bot ) return;
  if ( [ 'music', 'voting', 'notifications' ].indexOf( message.channel.name ) > -1 ) message.delete();

  if ( ( config.prefix.indexOf( message.content.charAt( 0 ) ) === -1 || message.content.length === 1 ) && !memberProcess[ message.author.id ] ) return;
  else if ( memberProcess[ message.author.id ] ) {
    try {
      return client.commands.get( 'member' ).execute( message, message.content );
    } catch ( error ) {
      message.reply( `ett fel uppstog med kommandot ${message.content}.` );
      return console.error( error );
    }
  }

  if ( !message.deleted && message.channel.type === 'text' ) message.delete( 10000 );

  // gets the args of the command
  let args = message.content.slice( 1 ).split( ' ' );
  args = args.filter( n => n );

  // gets the command
  const commandName = args.shift().toLowerCase();

  const command = client.commands.get( commandName ) || client.commands.find( cmd => cmd.aliases && cmd.aliases.includes( commandName ) );

  if ( !command ) {
    return message.reply( 'det finns inget sådant kommando.' ).then( ( msg ) => msg.delete( 10000 ) );
  }

  if ( command.serverOnly && message.channel.type !== 'text' ) {
    return message.reply( 'detta kommandot finns bara tillgänglig i Digital Ungdom servern.' );
  }

  if ( command.adminOnly ) {
    if ( !message.member.roles.find( r => r.name === 'admin' ) ) return message.reply( `du har inte behörighet att köra \`${command.name}\`.` ).then( ( msg ) => msg.delete( 10000 ) );
  }

  if ( command.group === 'music' ) {
    if ( message.channel.name !== 'music' ) {
      return message.reply( ' du måste skriva i `music` kanalen.' ).then( msg => { msg.delete( 10000 ); } );
    }
    if ( !message.member.voiceChannel ) {
      return message.reply( 'du måste vara i en ljud kanal.' ).then( msg => { msg.delete( 10000 ); } );
    }
    if ( guild.voiceConnection && ( message.member.voiceChannelID !== guild.me.voiceChannelID ) ) {
      return message.reply( 'du och boten måste vara i samma ljudkanal' ).then( msg => { msg.delete( 10000 ); } );
    }
  }

  try {
    command.execute( message, args );
  } catch ( error ) {
    message.reply( `Ett fel uppstog med kommandot ${message.content}.` );
    console.error( error );
  }
};
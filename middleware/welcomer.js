/* global guild include */

const Attachment = require( 'discord.js' ).Attachment;
const Canvas = require( 'canvas' );
const GIFEncoder = require( 'gifencoder' );
const fetch = require( 'node-fetch' );

const createNotificationEmbed = include( 'utils/embeds/createNotificationEmbed' );

module.exports = async function ( member ) {
  return new Promise( async ( resolve ) => {
    const channel = guild.channels.find( ch => ch.name === 'general' );
    if ( !channel ) return;

    const encoder = new GIFEncoder( 800, 250 );
    encoder.start();
    encoder.setRepeat( 0 ); // 0 for repeat, -1 for no-repeat
    encoder.setDelay( 150 );
    encoder.setQuality( 100 );

    const canvas = Canvas.createCanvas( 800, 250 );
    const ctx = canvas.getContext( '2d' );

    const buffer = await ( await fetch( member.user.displayAvatarURL ) ).buffer();
    const avatar = await Canvas.loadImage( buffer );

    ctx.fillStyle = '#3E5ABE';
    ctx.fillRect( 0, 0, canvas.width, canvas.height );

    ctx.strokeStyle = '#3E5ABE';
    ctx.strokeRect( 0, 0, canvas.width, canvas.height );

    let nameFont = 70;
    do {
      nameFont -= 1;
      ctx.font = `${ nameFont }pt Unifont`;
    } while ( ctx.measureText( `>${member.displayName}_` ).width > canvas.width - 275 );

    for ( let i = 0; i < member.displayName.length; i++ ) {
      ctx.fillStyle = '#ffffff';

      ctx.font = `${ nameFont }pt Unifont`;
      ctx.fillText( `>${member.displayName.slice( 0, i + 1 )}_`, 250, canvas.height / 2 + nameFont / 2 );

      ctx.drawImage( avatar, 25, 25, 200, 200 );
      encoder.addFrame( ctx );

      ctx.fillStyle = '#3E5ABE';
      ctx.fillRect( 0, 0, canvas.width, canvas.height );
    }

    for ( let i = 0; i < 25; i++ ) {
      ctx.fillStyle = '#ffffff';

      ctx.font = `${ nameFont }pt Unifont`;
      if ( i % 2 === 0 ) {
        ctx.fillText( `>${member.displayName}`, 250, canvas.height / 2 + nameFont / 2 );
      } else {
        ctx.fillText( `>${member.displayName}_`, 250, canvas.height / 2 + nameFont / 2 );
      }

      ctx.drawImage( avatar, 25, 25, 200, 200 );

      encoder.addFrame( ctx );
      encoder.addFrame( ctx );
      encoder.addFrame( ctx );

      ctx.fillStyle = '#3E5ABE';
      ctx.fillRect( 0, 0, canvas.width, canvas.height );
    }

    ctx.fillStyle = '#ffffff';

    ctx.font = `${ nameFont }pt Unifont`;
    ctx.fillText( `>${member.displayName}_`, 250, canvas.height / 2 + nameFont / 2 );

    ctx.drawImage( avatar, 25, 25, 200, 200 );

    encoder.finish();

    const gifBuffer = Buffer.from( encoder.out.data );

    const attachment = new Attachment( gifBuffer, 'välkommen.gif' );

    const notification = createNotificationEmbed( 'NEW MEMBER', 'welcome them', 4086462, { 'id': member.id, 'name': member.displayName, 'url': member.user.displayAvatarURL } );
    const notificationChannel = guild.channels.find( ch => ch.name === 'notifications' );
    notificationChannel.send( '@here, ny notifikation', { 'embed': notification } );

    channel.send( `Välkommen till servern, ${member}!`, attachment );

    resolve();
  } );
};
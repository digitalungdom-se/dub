/* global */

const ytdl = require( 'ytdl-core' );

module.exports = class Controller {
  constructor( client, guild ) {
    console.log( global.status );
    this.client = client;
    this.guild = guild;

    this.player;
    this.queue = [];
    this.volume = 0.1;
    this.connection = guild.VoiceConnection;

    this.musicChannel = guild.channels.find( ch => ch.name === 'music' );

    this.embed = {
      'title': 'Inget spelas',
      'color': 4086462
    };

    this.message = this.newController( { 'embed': this.embed } );
  }

  async play() {
    const url = this.queue[ 0 ].url;
    this.player = this.connection.playStream( ytdl( url, { quality: 'highestaudio', filter: 'audioonly' } ), {
      bitrate: 192000,
      volume: this.volume,
      passes: 3
    } );

    const metadata = await ytdl.getBasicInfo( url );
    let seconds = metadata.length_seconds % 60;
    if ( seconds < 10 ) seconds = `0${seconds}`;

    this.updateDisplay();

    this.client.user.setActivity( metadata.title, { type: 'LISTENING' } );

    if ( !this.message ) this.message = await this.musicChannel.send( { 'embed': this.embed } );

    this.player.on( 'end', function () {
      this.queue.shift();
      if ( this.queue.length !== 0 ) {
        this.play();
      } else {
        this.embed = {
          'title': 'Inget spelas',
          'color': 4086462
        };

        this.updateDisplay();
        this.queue = [];
        this.volume = 0.1;

        this.searchList = false;
        this.searchMessage = false;

        this.connection.disconnect();
        this.connection = false;

        this.client.user.setActivity( status.acitivity, { 'type': status.type } );
      }
    }.bind( this ) );
  }

  async add( song, message ) {
    this.queue.push( song );
    if ( !this.connection ) {
      if ( message.member ) this.connection = await message.member.voiceChannel.join();
      else this.connection = await message.voiceChannel.join();
      this.play();
    } else {
      this.updateDisplay();
    }
  }

  skip() {
    this.player.end();
  }

  stop() {
    this.queue = [];
    this.player.destroy();
  }

  setVolume( volume ) {
    if ( volume.set ) this.volume = parseFloat( volume.set );
    else if ( volume.inc ) this.volume += parseFloat( volume.inc );

    this.player.setVolume( this.volume );

    this.updateDisplay();
  }

  pauseResume() {
    if ( this.player.paused ) this.player.resume();
    else if ( !this.player.paused ) this.player.pause();
    this.updateDisplay();
  }

  updateDisplay() {
    if ( this.message ) {
      if ( this.queue[ 0 ] ) {
        this.embed = this.queue[ 0 ].embed;

        const volumeField = {
          'name': '🔊',
          'value': `${this.volume.toFixed(2)}`,
          'inline': true
        };

        const pauseField = {
          'name': '⏸',
          'value': `${this.player.paused}`,
          'inline': true
        };

        this.embed.fields = [ volumeField, pauseField ];

        if ( this.queue.length > 1 ) {
          const queueField = {};
          queueField.name = '__**Kö**__';
          queueField.value = '';
          for ( const [ index, song ] of this.queue.slice( 1 ).entries() ) {
            queueField.value += `**${index}.** ${song.metadata.author.name} | ${song.metadata.title}\n`;
          }

          this.embed.fields.push( queueField );
        }
      }
    }

    this.message.edit( { 'embed': this.embed } );
  }

  async newController() {
    this.message = await this.musicChannel.send( { 'embed': this.embed } );

    const reactions = [ '❎', '⏯', '⏭', '➕', '➖' ];
    for ( const reaction of reactions ) {
      await this.message.react( reaction );
    }
  }
};
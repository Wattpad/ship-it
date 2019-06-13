import React from 'react';
import ReactExpandableGrid from './components/TileGrid';
import './App.css';
import TopBar from './components/TopBar';

var data = [ // old image data (this json will be restructured with kube deploy information)
  { 'img': 'http://i.imgur.com/zIEjP6Q.jpg', 'link': 'https://www.instagram.com/p/BRFjVZtgSJD/', 'title': 'Westland Tai Poutini National Park', 'description': 'Photo by @christopheviseux / The Westland Tai Poutini National Park in New Zealand’s South Island offers a remarkable opportunity to take a guided walk on a glacier. A helicopter drop high on the Franz Josef Glacier, provides access to explore stunning ice formations and blue ice caves. Follow me for more images around the world @christopheviseux #newzealand #mountain #ice' },
  { 'img': 'http://i.imgur.com/rCrvQTv.jpg', 'link': 'https://www.instagram.com/p/BQ6_Wa2gmdR/', 'title': 'Dubai Desert Conservation Reserve', 'description': 'Photo by @christopheviseux / Early morning flight on a hot air balloon ride above the Dubai Desert Conservation Reserve. Merely an hour drive from the city, the park was created to protect indigenous species and biodiversity. The Arabian Oryx, which was close to extinction, now has a population well over 100. There are many options to explore the desert and flying above may be one of the most mesmerizing ways. Follow me @christopheviseux for more images from the Middle East. #dubai #desert' },
  { 'img': 'http://i.imgur.com/U8iVzVl.jpg', 'link': 'https://www.instagram.com/p/BQyfDiKAEq9/', 'title': 'Crumbling Reflections', 'description': 'Photo @pedromcbride // Crumbling Reflections: Much has changed in Cuba over the 17 years I have visited this island. But much has stayed the same. Time still ticks at a Cuban pace and old cars still run… I don’t know how... and while pockets of new construction and renovation exist thanks to a growing tourism boom, most buildings are crumbling and cracking under the Caribbean climate. But amidst the hardship, nostalgia and messy vitality, the Cuban people keep moving, like their cars. And somehow, they do it with a colorful friendliness and warmth that always amazes me. To see more, follow @pedromcbride #cuba #havana #photo #workshop @natgeoexpeditions #reflection #photooftheday #petemcbride.' },
  { 'img': 'http://i.imgur.com/Ky9aJlE.jpg', 'link': 'https://www.instagram.com/p/BQxf6CEgD8p/', 'title': 'Impalas', 'description': 'Impetious young impala go head-to-head as they practice sparring. A talent they will need later in life when the rut begins. Photographed on assignment for @natgeotravel in Kruger National Park. For more images from Kruger, South Africa, follow @kengeiger #natgeotravel #krugernationalpark' },
  { 'img': 'http://i.imgur.com/mf3qfzt.jpg', 'link': 'https://www.instagram.com/p/BQvy7gbgynF/', 'title': 'Elephants', 'description': 'Photo by @ronan_donovan // Two bull African elephants at dawn in Uganda\'s Murchison Falls National Park. See more from Uganda with @ronan_donovan.' },
  { 'img': 'http://i.imgur.com/zIEjP6Q.jpg', 'link': 'https://www.instagram.com/p/BRFjVZtgSJD/', 'title': 'Westland Tai Poutini National Park', 'description': 'Photo by @christopheviseux / The Westland Tai Poutini National Park in New Zealand’s South Island offers a remarkable opportunity to take a guided walk on a glacier. A helicopter drop high on the Franz Josef Glacier, provides access to explore stunning ice formations and blue ice caves. Follow me for more images around the world @christopheviseux #newzealand #mountain #ice' },
  { 'img': 'http://i.imgur.com/rCrvQTv.jpg', 'link': 'https://www.instagram.com/p/BQ6_Wa2gmdR/', 'title': 'Dubai Desert Conservation Reserve', 'description': 'Photo by @christopheviseux / Early morning flight on a hot air balloon ride above the Dubai Desert Conservation Reserve. Merely an hour drive from the city, the park was created to protect indigenous species and biodiversity. The Arabian Oryx, which was close to extinction, now has a population well over 100. There are many options to explore the desert and flying above may be one of the most mesmerizing ways. Follow me @christopheviseux for more images from the Middle East. #dubai #desert' },
  { 'img': 'http://i.imgur.com/U8iVzVl.jpg', 'link': 'https://www.instagram.com/p/BQyfDiKAEq9/', 'title': 'Crumbling Reflections', 'description': 'Photo @pedromcbride // Crumbling Reflections: Much has changed in Cuba over the 17 years I have visited this island. But much has stayed the same. Time still ticks at a Cuban pace and old cars still run… I don’t know how... and while pockets of new construction and renovation exist thanks to a growing tourism boom, most buildings are crumbling and cracking under the Caribbean climate. But amidst the hardship, nostalgia and messy vitality, the Cuban people keep moving, like their cars. And somehow, they do it with a colorful friendliness and warmth that always amazes me. To see more, follow @pedromcbride #cuba #havana #photo #workshop @natgeoexpeditions #reflection #photooftheday #petemcbride.' },
  { 'img': 'http://i.imgur.com/Ky9aJlE.jpg', 'link': 'https://www.instagram.com/p/BQxf6CEgD8p/', 'title': 'Impalas', 'description': 'Impetious young impala go head-to-head as they practice sparring. A talent they will need later in life when the rut begins. Photographed on assignment for @natgeotravel in Kruger National Park. For more images from Kruger, South Africa, follow @kengeiger #natgeotravel #krugernationalpark' },
  { 'img': 'http://i.imgur.com/mf3qfzt.jpg', 'link': 'https://www.instagram.com/p/BQvy7gbgynF/', 'title': 'Elephants', 'description': 'Photo by @ronan_donovan // Two bull African elephants at dawn in Uganda\'s Murchison Falls National Park. See more from Uganda with @ronan_donovan.' },
  { 'img': 'http://i.imgur.com/zIEjP6Q.jpg', 'link': 'https://www.instagram.com/p/BRFjVZtgSJD/', 'title': 'Westland Tai Poutini National Park', 'description': 'Photo by @christopheviseux / The Westland Tai Poutini National Park in New Zealand’s South Island offers a remarkable opportunity to take a guided walk on a glacier. A helicopter drop high on the Franz Josef Glacier, provides access to explore stunning ice formations and blue ice caves. Follow me for more images around the world @christopheviseux #newzealand #mountain #ice' },
  { 'img': 'http://i.imgur.com/rCrvQTv.jpg', 'link': 'https://www.instagram.com/p/BQ6_Wa2gmdR/', 'title': 'Dubai Desert Conservation Reserve', 'description': 'Photo by @christopheviseux / Early morning flight on a hot air balloon ride above the Dubai Desert Conservation Reserve. Merely an hour drive from the city, the park was created to protect indigenous species and biodiversity. The Arabian Oryx, which was close to extinction, now has a population well over 100. There are many options to explore the desert and flying above may be one of the most mesmerizing ways. Follow me @christopheviseux for more images from the Middle East. #dubai #desert' },
  { 'img': 'http://i.imgur.com/rCrvQTv.jpg', 'link': 'https://www.instagram.com/p/BQ6_Wa2gmdR/', 'title': 'Dubai Desert Conservation Reserve', 'description': 'Photo by @christopheviseux / Early morning flight on a hot air balloon ride above the Dubai Desert Conservation Reserve. Merely an hour drive from the city, the park was created to protect indigenous species and biodiversity. The Arabian Oryx, which was close to extinction, now has a population well over 100. There are many options to explore the desert and flying above may be one of the most mesmerizing ways. Follow me @christopheviseux for more images from the Middle East. #dubai #desert' }
]

var newData = [
  {
    serviceName: "loki",
    deploymentStatus: "deployed",
    timeStamp: "5 min ago",
    chartVersion: "0.1.1",
    autoDeploy: true,
    mirandaPR: "https://github.com",
    highlanderPR: "https://github.com",
    slackLink: "https://slack.com",
    dockerLink: "https://www.docker.com/",
    dataDogDashboard: "dd",
    dataDogMonitor: "",
    sumoLogs: "",
    travisBuild: "",
    kubeResources: [
      "pod 1",
      "pod 2",
      "pod 3",
      "pod 4",
    ]
  }, {
    serviceName: "hermes",
    deploymentStatus: "deployed",
    timeStamp: "5 min ago",
    chartVersion: "0.1.1",
    autoDeploy: true,
    mirandaPR: "https://github.com",
    highlanderPR: "https://github.com",
    slackLink: "https://slack.com",
    dockerLink: "https://www.docker.com/",
    dataDogDashboard: "dd",
    dataDogMonitor: "",
    sumoLogs: "",
    travisBuild: "",
    kubeResources: [
      "pod 1",
      "pod 2",
      "pod 3",
      "pod 4",
    ]
  },
]

function App() {
  var dataString = JSON.stringify(data)
  return (
    <div className="App">
      <TopBar />
      <ReactExpandableGrid
        gridData={dataString}
        detailHeight={300}
        ExpandedDetail_image_size={300}
        cellSize={250}
        ExpandedDetail_closeX_bool={false}
      />
    </div>
  );
}

export default App;

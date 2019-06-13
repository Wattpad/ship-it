import React from 'react'

import HelmIcon from '../assets/helm_icon.png'
import KubeIcon from '../assets/kubernetes_icon.png'
import DataDogIcon from '../assets/data_dog_icon.png'
import TravisIcon from '../assets/travis_icon.png'
import SumoIcon from '../assets/sumo_logic_icon.png'
import Typography from '@material-ui/core/Typography'
import Switch from '@material-ui/core/Switch'
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import { Paper, ListItem, ListItemText, List } from '@material-ui/core';

const theme = createMuiTheme({
  palette: {
    primary: {
      main: '#FF6612'
    },
    secondary: {
      main: '#FEAF0A'
      //main: '#FFFFFF'
    }
  }
})

const imgAlt = "not found"

class ExpandedDetail extends React.Component {
  constructor(props) {
    super(props)
    this.state = {}
  }

  render() {
    return (
      <MuiThemeProvider theme={theme}>
        <div className="flex-container">
          <div className="helm-status">
            <div className="right-padded">
              <img src={HelmIcon} width="32" height="32" alt={imgAlt} />
            </div>
            <div>
              <Typography variant="h5">ChartVersion</Typography>
            </div>
          </div>

          <div className="switch-status">
            <FormControlLabel
              control={
                <Switch color="primary" />
              }
              label="AutoDeploy"
            />
          </div>

          <div className="dataDog-status">
            <div className="right-padded">
              <img src={DataDogIcon} width="32" height="32" alt={imgAlt} />
            </div>
            <div>
              <Typography variant="h6">Dashboard | Monitor</Typography>
            </div>

          </div>
        </div>
        <div className="flex-container">
          <div className="kube-status">
            <div className="right-padded">
              <img src={KubeIcon} width="32" height="32" alt={imgAlt} />
            </div>
            <div>
              <Typography variant="h6">Resources</Typography>
            </div>
          </div>
          <div className="log-status">
            <div className="right-padded">
              <img src={SumoIcon} width="32" height="32" alt={imgAlt} />
            </div>
            <div>
              <Typography variant="h6">Logs</Typography>
            </div>
            <div className="double-padded">
              <img src={TravisIcon} width="32" height="32" alt={imgAlt} />
            </div>
            <div>
              <Typography variant="h6">Build</Typography>
            </div>
          </div>
        </div>
        <div className="kube-resource-list">
          <Paper style={{ maxHeight: '200px', minWidth: '100%', overflow: 'auto' }}>
            <List
              component="nav"
            >
              <ListItem button>
                <ListItemText>Pod</ListItemText>
              </ListItem>
              <ListItem button>
                <ListItemText>Pod</ListItemText>
              </ListItem>
              <ListItem button>
                <ListItemText>Pod</ListItemText>
              </ListItem>
              <ListItem button>
                <ListItemText>Pod</ListItemText>
              </ListItem>
              <ListItem button>
                <ListItemText>Pod</ListItemText>
              </ListItem>
              <ListItem button>
                <ListItemText>Pod</ListItemText>
              </ListItem>
              <ListItem button>
                <ListItemText>Pod</ListItemText>
              </ListItem>
            </List>
          </Paper>
        </div>
      </MuiThemeProvider>
    )
  }
}

export default ExpandedDetail;
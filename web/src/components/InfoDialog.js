import React from 'react'
import Dialog from '@material-ui/core/Dialog'
import { DialogContent, DialogActions, Button, DialogTitle, Typography } from '@material-ui/core';

class InfoDialog extends React.Component {
  render() {
    return (
      <Dialog open={this.props.open} onClose={this.props.handleClose}>
        <DialogTitle>Service Ownership</DialogTitle>
        <DialogContent>
          <Typography>Owned By: {this.props.owner}</Typography>
          <Typography>Slack Channel: {this.props.slack}</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.props.handleClose}>Close</Button>
        </DialogActions>
      </Dialog>
    )
  }
}

export default InfoDialog
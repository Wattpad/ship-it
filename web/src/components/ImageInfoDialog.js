import React from 'react'
import Dialog from '@material-ui/core/Dialog'
import { DialogContent, DialogActions, List, Button, DialogTitle, Typography, ListItem, ListItemText } from '@material-ui/core';

class ImageInfoDialog extends React.Component {
    getRegistry(docker) {
        let arr = docker.image.split('/')
        return arr[0]
    }

    getRepo(docker) {
        let arr = docker.image.split('/')
        return arr[1]
    }

    getTag(docker) {
        return docker.tag
    }

    getURI(docker) {
        return docker.image + ':' + this.getTag()
    }

    render() {
        this.images = this.props.docker.map(d => 
            <div>
                <ListItem>
                    <ListItemText>
                        <Typography><b>Registry:</b> {this.getRegistry(d)}</Typography>
                        <Typography><b>Repository:</b> {this.getRepo(d)}</Typography>
                        <Typography><b>Tag:</b> {this.getTag(d)}</Typography>
                    </ListItemText>
                </ListItem>
            </div>
        )
        
        return (
            <Dialog open={this.props.open} onClose={this.props.handleClose} maxWidth="xl">
                <DialogTitle>Docker Image Information</DialogTitle>
                <DialogContent>
                    <List>
                        {this.images}
                    </List>
                </DialogContent>
                <DialogActions>
                    <Button onClick={this.props.handleClose}>Close</Button>
                </DialogActions>
            </Dialog>
        )
    }
}

export default ImageInfoDialog;

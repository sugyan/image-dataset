import React from "react";
import { Link as RouterLink, LinkProps } from "react-router-dom";
import { Link, List, ListItem } from "@material-ui/core";

const Menu: React.FC = () => {
    const link = React.forwardRef<HTMLAnchorElement, LinkProps>((props, ref) => {
        return <RouterLink innerRef={ref} {...props} />;
    });
    return (
      <div>
        <List>
          <ListItem>
            <Link component={link} to="/images">Images</Link>
          </ListItem>
          <ListItem>
            <Link component={link} to="/stats">Stats</Link>
          </ListItem>
        </List>
      </div>
    );
};

export default Menu;

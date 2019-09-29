import React from "react";
import { Link as RouterLink, LinkProps } from "react-router-dom";
import { Link, List, ListItem } from "@material-ui/core";

const Menu: React.FC = () => {
    const link = React.forwardRef<HTMLAnchorElement, Omit<LinkProps, "to">>(
        (props, ref) => <RouterLink innerRef={ref} to="/images" {...props} />,
    );
    return (
      <div>
        <List>
          <ListItem>
            <Link component={link}>Images</Link>
          </ListItem>
        </List>
      </div>
    );
};

export default Menu;

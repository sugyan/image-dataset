import CSS from "csstype";
import React, { useEffect, useState } from "react";
import { useHistory, useLocation } from "react-router";
import { Link as RouterLink, LinkProps } from "react-router-dom";
import {
    Card, CardActionArea, CardMedia, CardContent,
    Box, CircularProgress, Grid, Link, Typography,
    Theme, makeStyles, createStyles,
} from "@material-ui/core";

import SearchBox from "./SearchBox";
import { ImageResponse } from "../common/interfaces";

const useStyles = makeStyles((theme: Theme) => {
    return createStyles({
        progress: {
            margin: theme.spacing(2),
        },
        card: {
            maxWidth: 144,
        },
        cardContent: {
            padding: theme.spacing(1.5),
        },
        media: {
            height: 144,
            width: 144,
        },
    });
});

const Images: React.FC = () => {
    const classes = useStyles();
    const history = useHistory();
    const location = useLocation();
    const [images, setImages] = useState<ImageResponse[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    useEffect(() => {
        if (location.search.length === 0) {
            return;
        }
        const params = new URLSearchParams(location.search);
        params.set("count", "100");
        setImages([]);
        setLoading(true);
        fetch(`/api/images?${params}`).then((res: Response) => {
            if (res.ok) {
                return res.json();
            }
            if (res.status === 401) {
                history.push("/");
                return;
            }
            throw new Error(res.statusText);
        }).then((data: ImageResponse[]) => {
            setImages(data);
        }).catch((err: Error) => {
            window.console.error(err.message);
        }).finally(() => {
            setLoading(false);
        });
    }, [history, location]);
    const cards = images.map((image: ImageResponse) => {
        const link = React.forwardRef<HTMLAnchorElement, Omit<LinkProps, "to">>(
            (props, ref) => {
                const to = {
                    pathname: `/image/${image.id}`,
                    search: location.search,
                };
                return (
                  <RouterLink innerRef={ref} to={to} {...props} />
                );
            },
        );
        const style: CSS.Properties = {
            borderStyle: "solid",
            borderWidth: "2px",
            borderColor: "transparent",
        };
        switch (image.status) {
        case 1:
            style.borderColor = "lightcoral";
            break;
        case 2:
            style.borderColor = "gray";
            break;
        case 3:
            style.borderColor = "lightgreen";
            break;
        }
        return (
          <Box key={image.id} m={0.2}>
            <Card className={classes.card} style={style}>
              <Link component={link}>
                <CardActionArea>
                  <CardMedia
                      className={classes.media}
                      image={image.image_url} />
                </CardActionArea>
              </Link>
              <CardContent className={classes.cardContent}>
                <Typography variant="body2">
                  {image.size}×{image.size}
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  {image.label_name}
                </Typography>
              </CardContent>
            </Card>
          </Box>
        );
    });
    const progress = (
      <Grid container justify="center">
        <CircularProgress className={classes.progress} />
      </Grid>
    );
    return (
      <React.Fragment>
        <h2>Images</h2>
        <SearchBox />
        {loading && progress}
        <Box display="flex" flexWrap="wrap" mt={2}>
          {cards}
        </Box>
      </React.Fragment>
    );
};

export default Images;

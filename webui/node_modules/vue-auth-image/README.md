# vue-auth-image

> A reusable directive for [Vue.js](https://github.com/vuejs/vue) that loads
> an image requiring authentication and includes it as data in-line in your web
> pages.

[![npm version](https://img.shields.io/npm/v/vue-auth-image.svg)](https://www.npmjs.com/package/vue-auth-image)

## Overview

Contrary to other HTTP requests, browsers don't send common headers such as
`Authorization` when retrieving an image specified in a `<img>` tag.

This Vue.js directive overcomes this limitation by providing a way to load your
images as any other resources and then embedding them into your web pages using
the `data:image/FILETYPE;base64` URI scheme.

## Requirements

- vue: \^2.0.0
- axios: >= 0.5.0

## Install

### npm

``` sh
$ npm install vue-auth-image --save
```

### Using a CDN

``` html
<script src="https://cdn.jsdelivr.net/npm/vue-auth-image/vue-auth-image.js"></script>
<!-- OR -->
<script src="https://cdn.jsdelivr.net/npm/vue-auth-image/vue-auth-image.min.js"></script>
```

### Using Nuxt.js

To use `vue-auth-image` with [Nuxt.js](https://nuxtjs.org/), start by creating
a new plugin file name `vue-auth-image.js` in your `plugins/` directory. Add
the following content to it:

```javascript
import Vue from 'vue';
import VueAuthImage from 'vue-auth-image';

Vue.use(VueAuthImage);
```

Then, add the plugin to your `nuxt.config.js` at the root of your project:

```javascript
plugins: [
  '@/plugins/vue-auth-image.js'
]
```

## API

### auth-image

A directive that requests an image URI asynchronously and embed it into your
`<img>` tag using the [data URI scheme](https://en.wikipedia.org/wiki/Data_URI_scheme).

``` js
import axios from 'axios';
import Vue from 'vue';
import VueAuthImage from 'vue-auth-image';

// register vue-auth-image directive
Vue.use(VueAuthImage);

// set Authorization header used by axios
var authHeader = 'Bearer ' + localStorage.getItem('id_token');
axios.defaults.headers.common['Authorization'] = authHeader;
```

Once the directive is registered, you can use it in your Vue templates.

``` html
<template>
  <div>
    <img v-auth-image="https://api.app.com/images/authenticatedImg.png">
  </div>
</template>
```

See `/example` for a demo. To build it, run `npm install && npm run build`.

## License

[MIT](https://opensource.org/licenses/MIT)

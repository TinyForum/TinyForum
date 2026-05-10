(function (e) {
  var t = Object.create,
    n = Object.defineProperty,
    r = Object.getOwnPropertyDescriptor,
    i = Object.getOwnPropertyNames,
    a = Object.getPrototypeOf,
    o = Object.prototype.hasOwnProperty,
    s = (e, t) => () => (
      t || (e((t = { exports: {} }).exports, t), (e = null)),
      t.exports
    ),
    c = (e, t, a, s) => {
      if ((t && typeof t == `object`) || typeof t == `function`)
        for (var c = i(t), l = 0, u = c.length, d; l < u; l++)
          ((d = c[l]),
            !o.call(e, d) &&
              d !== a &&
              n(e, d, {
                get: ((e) => t[e]).bind(null, d),
                enumerable: !(s = r(t, d)) || s.enumerable,
              }));
      return e;
    };
  e = ((e, r, i) => (
    (i = e == null ? {} : t(a(e))),
    c(
      r || !e || !e.__esModule
        ? n(i, `default`, { value: e, enumerable: !0 })
        : i,
      e,
    )
  ))(e, 1);
  var l = s((e) => {
      var t = Symbol.for(`react.transitional.element`);
      function n(e, n, r) {
        var i = null;
        if (
          (r !== void 0 && (i = `` + r),
          n.key !== void 0 && (i = `` + n.key),
          `key` in n)
        )
          for (var a in ((r = {}), n)) a !== `key` && (r[a] = n[a]);
        else r = n;
        return (
          (n = r.ref),
          {
            $$typeof: t,
            type: e,
            key: i,
            ref: n === void 0 ? null : n,
            props: r,
          }
        );
      }
      ((e.jsx = n), (e.jsxs = n));
    }),
    u = s((e, t) => {
      t.exports = l();
    })();
  function d({ postId: t }) {
    let [n, r] = (0, e.useState)(0);
    return (0, u.jsxs)(`div`, {
      className: `card bg-primary/10 border border-primary/20 p-3 text-sm`,
      children: [
        (0, u.jsx)(`p`, {
          className: `font-semibold text-primary`,
          children: `📌 我的插件`,
        }),
        t &&
          (0, u.jsxs)(`p`, {
            className: `text-xs text-base-content/50`,
            children: [`当前帖子：`, t],
          }),
        (0, u.jsxs)(`button`, {
          className: `btn btn-xs btn-primary mt-2`,
          onClick: () => r((e) => e + 1),
          children: [`点击了 `, n, ` 次`],
        }),
      ],
    });
  }
  var f = `demo`;
  window[`__plugin_${f}__`] = async function (t) {
    (t.getConfig().title,
      t.registerSlot(`sidebar-top`, () => e.default.createElement(d, {}), {
        order: 10,
      }),
      t.registerSlot(
        `post-detail-bottom`,
        ({ postId: t }) => e.default.createElement(d, { postId: t }),
        { order: 5 },
      ),
      t.on(`post:view`, (e) => {
        t.log(`info`, `用户浏览了帖子: ${JSON.stringify(e)}`);
      }),
      t.on(`user:login`, (e) => {
        let n = t.getUser();
        t.log(`info`, `用户 ${n?.username} 已登录`);
      }));
  };
})(React);

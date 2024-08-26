/**
 * Toasty is a class that provides a simple toast notification system for web applications.
 * It allows you to display customizable toast notifications on the screen, with various
 * configuration options such as position, duration, and styling.
 *
 * The class provides methods to push new toast notifications, close existing ones, and
 * manage the overall state of the toast stack.
 *
 * Strongly influenced by https://codeshack.io/elegant-toast-notifications-javascript/
 *
 * @class Toasty
 * @param {Object} [options] - The configuration options for the Toasty instance.
 * @param {string} [options.position="bottom-center"] - The position of the toast notifications on the screen.
 * @param {Array} [options.stack=[]] - The stack of active toast notifications.
 * @param {number} [options.offsetX=20] - The horizontal offset of the toast notifications from the screen edge.
 * @param {number} [options.offsetY=20] - The vertical offset of the toast notifications from the screen edge.
 * @param {number} [options.gap=20] - The gap between consecutive toast notifications.
 * @param {number} [options.numToasts=0] - The number of active toast notifications.
 * @param {string} [options.duration=".5s"] - The duration of the toast notification animation.
 * @param {string} [options.timing="ease"] - The timing function for the toast notification animation.
 * @param {boolean} [options.dimOld=true] - Whether to dim the old toast notifications when a new one is displayed.
 * @param {number} [options.zIndex=9999] - The z-index of the toast notifications.
 * @param {number} [options.autoCloseTime=3000] - The time in milliseconds before a toast notification is automatically closed.
 * @param {Object} [options.toastStyle] - The style of the toast notifications.
 * @param {Object} [options.toastStyle.success] - The style of success toast notifications.
 * @param {string} [options.toastStyle.success.main] - The main tailwind classes of success toast notifications.
 * @param {string} [options.toastStyle.success.header] - The header tailwind classes of success toast notifications.
 * @param {string} [options.toastStyle.success.body] - The body tailwind classes of success toast notifications.
 * @param {string} [options.toastStyle.success.closeIcon] - The header tailwind classes of success toast notifications.
 * @param {Object} [options.toastStyle.error] - The style of error toast notifications.
 * @param {string} [options.toastStyle.error.main] - The main tailwind classes of error toast notifications.
 * @param {string} [options.toastStyle.error.header] - The header tailwind classes of error toast notifications.
 * @param {string} [options.toastStyle.error.body] - The body tailwind classes of error toast notifications.
 * @param {string} [options.toastStyle.error.closeIcon] - The header tailwind classes of error toast notifications.
 * @param {Object} [options.toastStyle.warning] - The style of warning toast notifications.
 * @param {string} [options.toastStyle.warning.main] - The main tailwind classes of warning toast notifications.
 * @param {string} [options.toastStyle.warning.header] - The header tailwind classes of warning toast notifications.
 * @param {string} [options.toastStyle.warning.body] - The body tailwind classes of warning toast notifications.
 * @param {string} [options.toastStyle.warning.closeIcon] - The header tailwind classes of warning toast notifications.
 * @param {Object} [options.toastStyle.theme] - The style of theme toast notifications.
 * @param {string} [options.toastStyle.theme.main] - The main tailwind classes of theme toast notifications.
 * @param {string} [options.toastStyle.theme.header] - The header tailwind classes of theme toast notifications.
 * @param {string} [options.toastStyle.theme.body] - The body tailwind classes of theme toast notifications.
 * @param {string} [options.toastStyle.theme.closeIcon] - The header tailwind classes of theme toast notifications.
 * @param {Object} [options.toastStyle.info] - The style of info toast notifications.
 * @param {string} [options.toastStyle.info.main] - The main tailwind classes of info toast notifications.
 * @param {string} [options.toastStyle.info.header] - The header tailwind classes of info toast notifications.
 * @param {string} [options.toastStyle.info.body] - The body tailwind classes of info toast notifications.
 * @param {string} [options.toastStyle.info.closeIcon] - The header tailwind classes of info toast notifications.
 */
class Toasty {
  constructor(options) {
    let defaults = {
      position: "bottom-center",
      stack: [],
      offsetX: 20,
      offsetY: 20,
      gap: 20,
      numToasts: 0,
      duration: ".5s",
      timing: "ease",
      dimOld: true,
      zIndex: 9999,
      autoCloseTime: 3000,
      toastStyle: {
        success: {
          main: " bg-[#D1FADF] border-[1px] border-[#0EB869] ",
          header: "text-[#0EB869] uppercase text-xs font-semibold",
          body: "text-black",
          closeIcon: "font-medium text-[#0EB869] hover:font-bold",
        },
        error: {
          main: " bg-[#FEE4E2] border-[1px] border-[#F04339] ",
          header: "text-[#F04339] uppercase text-xs font-semibold",
          body: "text-black",
          closeIcon: "font-medium text-[#F04339] hover:font-bold",
        },
        warning: {
          main: " bg-[#FDEFC6] border-[1px] border-[#F79007] ",
          header: "text-[#F79007] uppercase text-xs font-semibold",
          body: "text-black",
          closeIcon: "font-medium text-[#F79007] hover:font-bold",
        },
        info: {
          main: " bg-[#CEDEFD] border-[1px] border-[#2369F6] ",
          header: "text-[#2369F6] uppercase text-xs font-semibold",
          body: "text-black",
          closeIcon: "font-medium text-[#2369F6] hover:font-bold",
        },
        theme: {
          main: " bg-background dark:bg-[#1F1F1F] border-[1px] ",
          header: "text-foreground uppercase text-xs font-semibold",
          body: "text-muted-foreground",
          closeIcon: "font-medium text-foreground hover:font-bold",
        },
      },
    };
    this.options = Object.assign(defaults, options);
  }

  /**
   * Pushes a new toast notification to the stack and displays it on the screen.
   *
   * @param {Object} obj - The configuration object for the toast notification.
   * @param {string} [obj.link] - The URL to link the toast notification to.
   * @param {string} [obj.linkTarget] - The target for the link, e.g. "_blank".
   * @param {string} [obj.style] - Additional CSS classes to apply to the toast notification.
   * @param {string} [obj.title] - The title of the toast notification.
   * @param {string} [obj.content] - The content of the toast notification.
   * @param {boolean} [obj.closeButton=true] - Whether to display a close button on the toast notification.
   * @param {number} [obj.width] - The width of the toast notification.
   * @param {function} [obj.onOpen] - A callback function to be called when the toast notification is opened.
   * @param {number} [obj.zIndex] - The z-index of the toast notification.
   */
  push(obj) {
    this.numToasts++;

    let toast = document.createElement(obj.link ? "a" : "div");
    let toastStyle = this.options.toastStyle[obj.style ?? "theme"];

    if (obj.link) {
      toast.href = obj.link;
      toast.target = obj.linkTarget ? obj.linkTarget : "_self";
    }

    toast.className =
      `fixed z-${
        "[" + (obj.zIndex ?? this.options.zIndex) + "]"
      } max-w-sm  shadow-lg rounded-lg flex px-4 py-2 transform -translate-y-full ` +
      toastStyle.main +
      (obj.style ? " " + obj.style : "") +
      " " +
      this.position;
    toast.innerHTML = `
            <div class="flex-1 pr-4 overflow-hidden">
                ${
                  obj.title
                    ? '<h3 class="mb-1 font-medium text-sm ' +
                      toastStyle.header +
                      ' break-words">' +
                      obj.title +
                      "</h3>"
                    : ""
                }
                ${
                  obj.content
                    ? '<div class="text-sm ' +
                      toastStyle.body +
                      ' break-words">' +
                      obj.content +
                      "</div>"
                    : ""
                }
            </div>
            ${
              obj.closeButton == null || obj.closeButton === true
                ? '<button class="appearance-none border-none bg-transparent cursor-pointer text-2xl leading-6 pb-1 ' +
                  toastStyle.closeIcon +
                  ' ">&times;</button>'
                : ""
            }
        `;
    document.body.appendChild(toast);
    toast.getBoundingClientRect();
    if (this.position == "top-left") {
      toast.style.top = 0;
      toast.style.left = this.offsetX + "px";
    } else if (this.position == "top-center") {
      toast.style.top = 0;
      toast.style.left = 0;
    } else if (this.position == "top-right") {
      toast.style.top = 0;
      toast.style.right = this.offsetX + "px";
    } else if (this.position == "bottom-left") {
      toast.style.bottom = 0;
      toast.style.left = this.offsetX + "px";
    } else if (this.position == "bottom-center") {
      toast.style.bottom = 0;
      toast.style.left = 0;
    } else if (this.position == "bottom-right") {
      toast.style.bottom = 0;
      toast.style.right = this.offsetX + "px";
    }
    if (obj.width || this.width) {
      toast.style.width = (obj.width || this.width) + "px";
    }
    toast.dataset.transitionState = "queue";
    let index = this.stack.push({
      element: toast,
      props: obj,
      offsetX: this.offsetX,
      offsetY: this.offsetY,
      index: 0,
    });
    this.stack[index - 1].index = index - 1;
    if (toast.querySelector("button")) {
      toast.querySelector("button").onclick = (event) => {
        event.preventDefault();
        this.closeToast(this.stack[index - 1]);
      };
    }
    if (obj.link) {
      toast.onclick = () => this.closeToast(this.stack[index - 1]);
    }
    this.openToast(this.stack[index - 1]);
    if (obj.onOpen) obj.onOpen(this.stack[index - 1]);

    // Auto close the toast after the specified time
    if (obj.autoCloseTime !== false) {
      const autoCloseTime = obj.autoCloseTime || this.options.autoCloseTime;
      setTimeout(() => {
        this.closeToast(this.stack[index - 1]);
      }, autoCloseTime);
    }
  }

  /**
   * Opens a toast notification and manages the transition and positioning of the toast.
   *
   * @param {Object} toast - An object containing the toast element and its properties.
   * @param {HTMLElement} toast.element - The toast element to be displayed.
   * @param {Object} toast.props - The properties of the toast, such as dismissAfter.
   * @param {number} toast.offsetX - The horizontal offset of the toast.
   * @param {number} toast.offsetY - The vertical offset of the toast.
   * @param {number} toast.index - The index of the toast in the stack.
   * @returns {boolean} - True if the toast was opened successfully, false otherwise.
   */
  openToast(toast) {
    if (this.isOpening() === true) {
      return false;
    }
    toast.element.dataset.transitionState = "opening";
    toast.element.style.transition =
      this.duration + " transform " + this.timing;
    this._transformToast(toast);
    toast.element.addEventListener("transitionend", () => {
      if (toast.element.dataset.transitionState == "opening") {
        toast.element.dataset.transitionState = "complete";
        for (let i = 0; i < this.stack.length; i++) {
          if (this.stack[i].element.dataset.transitionState == "queue") {
            this.openToast(this.stack[i]);
          }
        }
        if (toast.props.dismissAfter) {
          this.closeToast(toast, toast.props.dismissAfter);
        }
      }
    });
    for (let i = 0; i < this.stack.length; i++) {
      if (this.stack[i].element.dataset.transitionState == "complete") {
        this.stack[i].element.dataset.transitionState = "opening";
        this.stack[i].element.style.transition =
          this.duration +
          " transform " +
          this.timing +
          (this.dimOld ? ", " + this.duration + " opacity ease" : "");
        if (this.dimOld) {
          this.stack[i].element.classList.add("opacity-30");
        }
        this.stack[i].offsetY += toast.element.offsetHeight + this.gap;
        this._transformToast(this.stack[i]);
      }
    }
    return true;
  }

  /**
   * Closes a toast notification and manages the transition and positioning of the remaining toasts.
   *
   * @param {Object} toast - An object containing the toast element and its properties.
   * @param {HTMLElement} toast.element - The toast element to be closed.
   * @param {Object} toast.props - The properties of the toast, such as onClose callback.
   * @param {number} toast.index - The index of the toast in the stack.
   * @param {number} [delay] - An optional delay in milliseconds before the toast is closed.
   * @returns {boolean} - True if the toast was closed successfully, false otherwise.
   */
  closeToast(toast, delay = null) {
    if (this.isOpening() === true) {
      setTimeout(() => this.closeToast(toast, delay), 100);
      return false;
    }
    if (toast.element.dataset.transitionState == "close") {
      return true;
    }
    if (toast.element.querySelector(".close-button")) {
      toast.element.querySelector(".close-button").onclick = null;
    }
    toast.element.dataset.transitionState = "close";
    toast.element.style.transition =
      ".2s opacity ease" + (delay ? " " + delay : "");
    toast.element.classList.add("opacity-0");
    toast.element.addEventListener("transitionend", () => {
      if (toast.element.dataset.transitionState == "close") {
        let offsetHeight = toast.element.offsetHeight;
        if (toast.props.onClose) toast.props.onClose(toast);
        toast.element.remove();
        for (let i = 0; i < toast.index; i++) {
          this.stack[i].element.style.transition =
            this.duration + " transform " + this.timing;
          this.stack[i].offsetY -= offsetHeight + this.gap;
          this._transformToast(this.stack[i]);
        }
        let focusedToast = this.getFocusedToast();
        if (focusedToast) {
          focusedToast.element.classList.remove("opacity-30");
        }
      }
    });
    return true;
  }

  /**
   * Checks if any toast notifications are currently in the opening transition state.
   *
   * @returns {boolean} - True if any toast is currently opening, false otherwise.
   */
  isOpening() {
    let opening = false;
    for (let i = 0; i < this.stack.length; i++) {
      if (this.stack[i].element.dataset.transitionState == "opening") {
        opening = true;
      }
    }
    return opening;
  }

  /**
   * Retrieves the currently focused toast notification from the stack.
   *
   * @returns {Object|boolean} - The focused toast notification object, or false if no toast is currently focused.
   */
  getFocusedToast() {
    for (let i = 0; i < this.stack.length; i++) {
      if (this.stack[i].offsetY == this.offsetY) {
        return this.stack[i];
      }
    }
    return false;
  }

  /**
   * Transforms the position of a toast notification based on the configured position.
   *
   * @param {Object} toast - The toast notification object to transform.
   * @returns {void}
   */
  _transformToast(toast) {
    if (this.position == "top-center") {
      toast.element.style.transform = `translate(calc(50vw - 50%), ${toast.offsetY}px)`;
    } else if (this.position == "top-right" || this.position == "top-left") {
      toast.element.style.transform = `translate(0, ${toast.offsetY}px)`;
    } else if (this.position == "bottom-center") {
      toast.element.style.transform = `translate(calc(50vw - 50%), -${toast.offsetY}px)`;
    } else if (
      this.position == "bottom-left" ||
      this.position == "bottom-right"
    ) {
      toast.element.style.transform = `translate(0, -${toast.offsetY}px)`;
    }
  }

  set stack(value) {
    this.options.stack = value;
  }

  get stack() {
    return this.options.stack;
  }

  set position(value) {
    this.options.position = value;
  }

  get position() {
    return this.options.position;
  }

  set offsetX(value) {
    this.options.offsetX = value;
  }

  get offsetX() {
    return this.options.offsetX;
  }

  set offsetY(value) {
    this.options.offsetY = value;
  }

  get offsetY() {
    return this.options.offsetY;
  }

  set gap(value) {
    this.options.gap = value;
  }

  get gap() {
    return this.options.gap;
  }

  set numToasts(value) {
    this.options.numToasts = value;
  }

  get numToasts() {
    return this.options.numToasts;
  }

  set width(value) {
    this.options.width = value;
  }

  get width() {
    return this.options.width;
  }

  set duration(value) {
    this.options.duration = value;
  }

  get duration() {
    return this.options.duration;
  }

  set timing(value) {
    this.options.timing = value;
  }

  get timing() {
    return this.options.timing;
  }

  set dimOld(value) {
    this.options.dimOld = value;
  }

  get dimOld() {
    return this.options.dimOld;
  }
}

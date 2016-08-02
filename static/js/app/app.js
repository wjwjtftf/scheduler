
var App = {
    init:function(){


        var url = window.location;
        var element = $('.navbar-nav li a').filter(function() {
            return this.href == url || url.href.indexOf(this.href) == 0;
        }).parent().addClass('active').parent().parent().addClass('in').parent();
        if (element.is('li')) {
            element.addClass('active');
        }

    }
};

$(function(){

    App.init();
});
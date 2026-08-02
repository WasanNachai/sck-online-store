*** Settings ***
Library    SeleniumLibrary
Test Teardown   Close All Browsers

*** Variables ***
${URL}    http://localhost/product/list
${BROWSER}    headlesschrome
${REMOTE_HUB_URL}

*** Test Cases ***
ทดสอบ สั่งซื้อสินค้า Horses and Unicorns Set ผ่านทุกหน้า product/list product/[id] shopping cart checkout payment orders/completed
    เข้าสู่เว็บไซต์ และตรวจสอบว่า redirect มาที่    /auth/login    login-page
    เข้าสู่ระบบ    login-username-input    user_3    login-password-input    P@ssw0rd
    ตรวจสอบสินค้าในหน้า product list    product-card-name-3    Horses and Unicorns Set    product-card-price-3    834.95
    เลือกดูสินค้า    product-card-name-3
    ตรวจสอบรายละเอียดสินค้า    Horses and Unicorns Set    CoolKidZ    834.95    16
    เพิ่มสินค้าลงตะกร้าและมีสินค้าจำนวน    1
    ตรวจสอบข้อมูลสินค้าในตะกร้า และไปที่หน้า checkout    Horses and Unicorns Set    834.95    16    834.95
    ใส่ที่อยู่จัดส่งสินค้า
    ...    ณัฐพล    ศรีสมบัติ
    ...    43/8 หมู่บ้านเปี่ยมสุข ถนนลาดพร้าว ซอย 63    กรุงเทพมหานคร
    ...    เขตวังทองหลาง    วังทองหลาง
    ...    10310    0891234567
    เลือกวิธีจัดส่งสินค้าเป็น    kerry
    ตรวจสอบค่าจัดส่งสินค้าของ Kerry เท่ากันกับ 50.00 บาท    kerry    50.00
    เลือกช่องทางการชำระเงินแบบ VISA Credit Card    Nattapon Srisombat    4719700591590995    0226    752
    ตรวจสอบราคารวมในหน้า checkout    834.95    16 Points    884.95
    กดชำระเงินและตรวจสอบหน้า payment    884.95
    ยืนยัน OTP
    ตรวจสอบหน้า order completed    kerry    KR

*** Keywords ***
เข้าสู่เว็บไซต์ และตรวจสอบว่า redirect มาที่
    [Arguments]    ${target-url}    ${target-element-locator}
    Open Browser    url=${URL}    browser=${BROWSER}    remote_url=${REMOTE_HUB_URL}
    Wait Until Location Is Not    location=${URL}
    Location Should Contain    ${target-url}
    Page Should Contain Element    id:${target-element-locator}

เข้าสู่ระบบ
    [Arguments]    ${username-input-locator}    ${username}    ${password-input-locator}    ${password}
    Wait Until Element Is Visible    id:${username-input-locator}
    Input Text    id:${username-input-locator}    ${username}
    Input Password    id:${password-input-locator}    ${password}
    Click Button    id:login-btn
    Wait Until Location Is    ${URL}
    Wait Until Element Is Visible    product-list

ตรวจสอบสินค้าในหน้า product list
    [Arguments]    ${card-name-locator}    ${expected-name}    ${card-price-locator}    ${expected-price}
    Wait Until Element Is Visible    id:${card-name-locator}
    Element Text Should Be    id:${card-name-locator}    ${expected-name}
    Element Text Should Be    id:${card-price-locator}    ฿${expected-price}

เลือกดูสินค้า
    [Arguments]    ${card-name-locator}
    Click Element    id:${card-name-locator}

ตรวจสอบรายละเอียดสินค้า
    [Arguments]    ${product-name}    ${product-brand}    ${product-thb-price}    ${product-point}
    Wait Until Element Is Visible    id:product-detail-product-name
    Element Text Should Be    id:product-detail-product-name    ${product-name}
    Element Text Should Be    id:product-detail-brand    ${product-brand}
    Element Text Should Be    id:product-detail-price-thb    ฿${product-thb-price}
    Element Text Should Be    id:product-detail-point    ${product-point} Points

เพิ่มสินค้าลงตะกร้าและมีสินค้าจำนวน
    [Arguments]    ${cart-qty}
    Click Button    id:product-detail-add-to-cart-btn
    Wait Until Element Contains    id:header-menu-cart-badge    ${cart-qty}

ตรวจสอบข้อมูลสินค้าในตะกร้า และไปที่หน้า checkout
    [Arguments]    ${product-name}    ${product-price}    ${product-point}    ${subtotal-price}
    Click Button    id:header-menu-cart-btn
    Wait Until Element Is Visible    id:product-3-price
    Element Text Should Be    id:product-3-name    ${product-name}
    Element Text Should Be    id:product-3-price    ฿${product-price}
    Element Text Should Be    id:product-3-point    ${product-point} Points
    Element Attribute Value Should Be    id:product-3-quantity-input    value    1
    Element Text Should Be    id:shopping-cart-subtotal-price    ฿${subtotal-price}
    Click Element    id:shopping-cart-checkout-btn

ใส่ที่อยู่จัดส่งสินค้า
    [Arguments]    ${firstname}    ${lastname}
    ...    ${address}    ${province}    ${district}
    ...    ${subdistrict}    ${zipcode}    ${phone-number}
    Input Text    id:shipping-form-first-name-input    ${firstname}
    Input Text    id:shipping-form-last-name-input    ${lastname}
    Input Text    id:shipping-form-address-input    ${address}
    Select From List By Label    id:shipping-form-province-select    ${province}
    Select From List By Label    id:shipping-form-district-select    ${district}
    Select From List By Label    id:shipping-form-sub-district-select    ${subdistrict}
    Element Attribute Value Should Be    id:shipping-form-zipcode-input    value    ${zipcode}
    Input Text    id:shipping-form-mobile-input    ${phone-number}

เลือกวิธีจัดส่งสินค้าเป็น
    [Arguments]    ${method}
    &{DELIVERY_METHOD}    Create Dictionary
    ...    kerry=id:shipping-method-1-card
    ...    thai_post=id:shipping-method-2-card
    ...    lineman=id:shipping-method-3-card
    Click Element    ${DELIVERY_METHOD}[${method}]

ตรวจสอบค่าจัดส่งสินค้าของ Kerry เท่ากันกับ 50.00 บาท
    [Arguments]    ${method}    ${fee}
    &{DELIVERY_METHOD}    Create Dictionary
    ...    kerry=id:shipping-method-1-fee
    ...    thai_post=id:shipping-method-2-fee
    ...    lineman=id:shipping-method-3-fee
    Element Text Should Be    ${DELIVERY_METHOD}[${method}]    ฿${fee}

เลือกช่องทางการชำระเงินแบบ VISA Credit Card
    [Arguments]    ${credit-card-name}    ${credit-card-number}    ${credit-card-expired-date}    ${credit-card-cvv}
    Click Element    id:payment-credit-input
    Input Text    id:payment-credit-form-fullname-input    ${credit-card-name}
    Input Text    id:payment-credit-form-card-number-input    ${credit-card-number}
    Input Text    id:payment-credit-form-expiry-input    ${credit-card-expired-date}
    Input Text    id:payment-credit-form-cvv-input    ${credit-card-cvv}

ตรวจสอบราคารวมในหน้า checkout
    [Arguments]    ${subtotal-price}    ${receive-point}    ${total-price}
    Element Should Be Visible    id:order-summary-subtotal-price
    Element Text Should Be    id:order-summary-subtotal-price    ฿${subtotal-price}
    Element Text Should Be    id:order-summary-receive-point-price    ${receive-point}
    Element Text Should Be    id:order-summary-total-payment-price    ฿${total-price}

กดชำระเงินและตรวจสอบหน้า payment
    [Arguments]    ${expected-amount}
    Click Button    id:payment-now-btn
    Wait Until Location Contains    /payment
    Page Should Contain    ฿${expected-amount}

ยืนยัน OTP
    Wait Until Element Is Visible    id:otp-input
    Click Button    Request OTP
    Input Text    id:otp-input    124532
    Click Button    OK

ตรวจสอบหน้า order completed
    [Arguments]    ${shipping-method}    ${tracking-prefix}
    Wait Until Element Is Visible    id:order-success-order-id
    Location Should Contain    /orders/completed
    Element Should Be Visible    id:order-success-order-payment-date
    Element Text Should Be    id:order-success-shipping-method    ${shipping-method}
    ${tracking-id}=    Get Text    id:order-success-tracking-id
    Should Match Regexp    ${tracking-id}    ^${tracking-prefix}-\\d{7,9}$
